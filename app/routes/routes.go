package routes

import (
	"web-navbar/app/controller"
	"web-navbar/controller"

	"github.com/gofiber/fiber/v2"
)

func Routes(r *fiber.App) {
	api := r.Group("/api")

	api.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "TES",
		})
	})

	api.Post("/register", controller.RegisterAccount)
	api.Post("/login", controller.Login)
}
