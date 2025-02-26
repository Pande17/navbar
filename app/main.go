package main

import (
	"log"
	"web-navbar/app/config"
	"web-navbar/app/database"
	"web-navbar/app/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.InitEnv()

	database.ConnectMongoDB()
	// database.CreateCollectionsAndIndexes(database.MongoClient)

	app := fiber.New(fiber.Config{
		Network: "tcp",
	})
	routes.Routes(app)

	log.Fatal(app.Listen(":3000"))
}
