package controller

import (
	"context"
	"os"
	"time"
	"web-navbar/database"
	"web-navbar/model"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func RegisterAccount(c *fiber.Ctx) error {
	var accReq struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Kelamin  string `json:"kelamin"`
		Alamat   string `json:"alamat"`
	}

	if err := c.BodyParser(&accReq); err != nil {
		return BadRequest(c, "Form can not be empty!", "body parser")
	}

	if _, err := govalidator.ValidateStruct(&accReq); err != nil {
		return BadRequest(c, "Input not valid", "Register validator")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(accReq.Password), bcrypt.DefaultCost)
	if err != nil {
		return Conflict(c, "Can not hash the password", "hashing password register")
	}

	accountCollection := database.ConnectCollection("account")

	var existingAccount model.Account
	filter := bson.M{"id_number": accReq.IDNumber}

	err = accountCollection.FindOne(context.TODO(), filter).Decode(&existingAccount)
	if err == nil {
		return BadRequest(c, "This ID Number already used", "existing account")
	} else if err != mongo.ErrNoDocuments {
		return Conflict(c, "Server error! Try again later", "existing account")
	}

	account := model.Account{
		IDNumber: accReq.IDNumber,
		Username: accReq.Username,
		Password: string(hash),
		// Balance: ,
		Base: model.Base{
			ID:        primitive.NewObjectID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		},
	}

	_, err = accountCollection.InsertOne(context.TODO(), &account)
	if err != nil {
		return Conflict(c, "Can not Register now! Try again later..", "insert data")
	}

	return OK(c, "Account created successfully!", account)
}

// Function to login to admin account
func Login(c *fiber.Ctx) error {
	var adminReq struct {
		AdminName     string `json:"admin_name" valid:"required~Nama tidak boleh kosong!, stringlength(1|30)~Nama harus antara 1 hingga 30 karakter!"`
		AdminPassword string `json:"admin_password" valid:"required~Password tidak boleh kosong!"`
	}

	// Parse the request body
	if err := c.BodyParser(&adminReq); err != nil {
		return BadRequest(c, "Input tidak valid! Silakan periksa kembali.", "Gagal mem-parsing body login")
	}

	// New variable to store admin login data
	var admin model.AdminAccount
	adminCollection := database.GetCollection("adminAcc")

	// Find admin account with the same account name as input name
	err := adminCollection.FindOne(context.TODO(), bson.M{"admin_name": adminReq.AdminName}).Decode(&admin)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Unauthorized(c, "Nama atau Password salah! Silakan coba lagi.", "Gagal menemukan nama admin saat login")
		}
		return Conflict(c, "Kesalahan database! Silakan coba lagi.", "Gagal menemukan nama admin saat login")
	}

	// Check if DeletedAt field already has a value (account has been deleted)
	if admin.DeletedAt != nil && !admin.DeletedAt.IsZero() {
		return AlreadyDeleted(c, "Akun ini telah dihapus! Silakan hubungi admin.", "Periksa akun admin yang dihapus", admin.DeletedAt)
	}

	// Hashing the input password with bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(admin.AdminPassword), []byte(adminReq.AdminPassword))
	if err != nil {
		return Unauthorized(c, "Nama atau Password salah! Silakan coba lagi.", "Gagal memverifikasi password saat login")
	}

	// Generate token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": admin.ID.Hex(),
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Retrieve the "SECRET" environment variable
	secret := os.Getenv("SECRET")
	if secret == "" {
		return Conflict(c, "Kunci rahasia tidak diset! Silakan hubungi admin.", "Gagal mengambil kunci rahasia")
	}

	// Use the secret key to sign the token
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return BadRequest(c, "Login gagal! Silakan coba lagi.", "Gagal menggunakan kunci rahasia")
	}

	// Set a cookie for admin
	c.Cookie(&fiber.Cookie{
		Name:     "authToken",
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour * 30),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "None",
		Path:     "/",
	})

	// Return success
	return OK(c, "Login berhasil! Selamat datang.", tokenString)
}
