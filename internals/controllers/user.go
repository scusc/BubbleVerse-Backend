package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/scusc/Bubbleverse-Backend/internals/db"
	"github.com/scusc/Bubbleverse-Backend/internals/models"
)

func SyncUser(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)
	email := c.Locals("email").(string)
	name := c.Locals("name").(string)

	var user models.User

	// Check if user exists in the database
	// If not, create a new user record
	result := db.DB.First(&user, "firebase_uid = ?", uid)

	if result.Error != nil {
		user = models.User{
			FirebaseUID: uid,
			Email: email,
			Name: name,	
		}
		db.DB.Create(&user)
	}
	return c.JSON(user)
}