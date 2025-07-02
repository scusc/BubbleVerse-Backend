package main

import (
	"log"
	"github.com/gofiber/fiber/v2"
    "github.com/scusc/BubbleVerse-Backend/internals/middleware"
)

func main() {
	app := fiber.New()
	middleware.InitFireBase()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!!!")
	})
	log.Fatal(app.Listen(":8080"))
}