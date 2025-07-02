package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/scusc/Bubbleverse-Backend/internals/db"
	"github.com/scusc/Bubbleverse-Backend/internals/middleware"
	"github.com/scusc/Bubbleverse-Backend/internals/models"
	"github.com/scusc/Bubbleverse-Backend/internals/controllers"
	"log"
)

func main() {
	app := fiber.New()
	middleware.InitFireBase()
	db.InitDB()

	// AutoMigrate is used to automatically create or update database tables
	// to match the structure of your Go structs (models).
	// In this case, &models.User{} is passed, so it ensures the "users" table
	// exists and its columns match the fields of the User struct.
	// If you add, remove, or change fields in your User struct, AutoMigrate
	// will update the database schema accordingly.
	// This is useful during development to keep your database in sync with your code.
	// Note: AutoMigrate will NOT delete unused columns or tables to avoid data loss.
	db.DB.AutoMigrate(&models.User{})


	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!!!")
	})

	// Firebase authentication middleware is applied to the /profile route.
	// This means that only authenticated users can access this protected route.
	app.Get("/profile", middleware.FirebaseAuthMiddleware, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"uid":   c.Locals("uid"),
			"email": c.Locals("email"),
			"name": c.Locals("name"),
		})
	})

	// The /me route is protected by the FirebaseAuthMiddleware.
	// When a user accesses this route, their Firebase authentication token is verified,
	// and if valid, the SyncUser controller is called.
	// This controller will either retrieve the user's information from the database
	// or create a new user record if it doesn't exist.
	app.Get("/me", middleware.FirebaseAuthMiddleware, controllers.SyncUser)

	log.Fatal(app.Listen(":8080"))
}