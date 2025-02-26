package model

type Accounts struct {
	Username string `json:"username" bson:"username"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	Kelamin  string `json:"kelamin" bson:"kelamin"`
	Alamat   string `json:"alamat" bson:"alamat"`
	Base     `bson:",inline"`
}
