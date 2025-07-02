package main

import (
	"log"
	"github.com/gofiber/fiber/v2"
    "github.com/scusc/Bubbleverse-Backend/internals/middleware"
)

func main() {
	app := fiber.New()
	middleware.InitFireBase()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!!!")
	})

	app.Get("/profile", middleware.FirebaseAuthMiddleware, func(c *fiber.Ctx) error {
		uid := c.Locals("uid")
		email := c.Locals("email")
		return c.JSON(fiber.Map{
			"uid":   uid,
			"email": email,
		})
	})

	log.Fatal(app.Listen(":8080"))
}