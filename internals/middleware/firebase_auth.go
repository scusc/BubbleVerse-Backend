package middleware

import (
	"context"
	"log"
	"strings"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/option"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
    secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

var firebaseAuthClient *auth.Client

func InitFireBase() {
	projectID := "71211314983" // Replace with your actual project ID
	secretID := "BubbleVerseFirebaseJson" // Replace with your secret ID
	versionID := "latest" // Or a specific version number like "1"

	ctx := context.Background()

	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create secretmanager client: %v", err)
	}
	defer client.Close()

	// Build the resource name of the secret version.
	name := fmt.Sprintf("projects/%s/secrets/%s/versions/%s", projectID, secretID, versionID)

	// Access the secret version.
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	resp, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		log.Fatalf("Failed to access secret version: %v", err)
	}

	tmpFile, err := os.CreateTemp("", "firebase-credentials-*.json")
	if err != nil {
		log.Fatalf("Failed to create temp file for credentials: %v", err)
	}
	defer os.Remove(tmpFile.Name()) 

	if _, err := tmpFile.Write(resp.Payload.Data); err != nil {
        log.Fatalf("Failed to write credentials to temp file: %v", err)
    }
    if err := tmpFile.Close(); err != nil {
        log.Fatalf("Failed to close temp file: %v", err)
    }

	opt := option.WithCredentialsFile(tmpFile.Name())
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

