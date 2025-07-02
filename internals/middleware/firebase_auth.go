package middleware

import (
	"context"
	"log"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/option"
)

var firebaseAuthClient *auth.Client

func InitFireBase() {
	opt := option.WithCredentialsFile("../../bubbleverse-94043-firebase-adminsdk-fbsvc-0fd5875ef7.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	authClient, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	firebaseAuthClient = authClient
	log.Println("Firebase Auth initialized")
}

func FirebaseAuthMiddleware(c *fiber.Ctx) error { 
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization header is missing",
		})
	}

	idToken := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := firebaseAuthClient.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		log.Printf("error verifying ID token: %v\n", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	c.Locals("uid", token.UID)
	c.Locals("email", token.Claims["email"])

	return c.Next()
}

