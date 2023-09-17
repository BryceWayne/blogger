package routes

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/BryceWayne/blogger/models"
	"github.com/BryceWayne/blogger/utils"
	"github.com/gofiber/fiber/v2"
)

func NewGitHubEvent(c *fiber.Ctx, config *utils.Config, client *firestore.Client) error {
	log.Println("INFO: New GitHub Event.")

	payload_ := c.Body()
	signature := c.Get("X-Hub-Signature")
	log.Printf("DEBUG: signature - %v", signature)

	secret := []byte(config.WebhookSecret) // Replace with your GitHub webhook secret

	if !verifySignature(secret, payload_, signature) {
		return c.Status(401).SendString("ERROR: Mismatched signature")
	}

	log.Println("DEBUG: Received valid payload")

	ctx := context.Background()
	var payload models.WebhookPayload
	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	log.Printf("DEBUG: Payload %+v\n", payload)

	// Loop through commits to log added, modified, and removed files
	for _, commit := range payload.Commits {
		log.Println("DEBUG: Added files:", commit.Added)
		log.Println("DEBUG: Modified files:", commit.Modified)
		log.Println("DEBUG: Removed files:", commit.Removed)
	}

	_, _, err := client.Collection("githubEvents").Add(ctx, payload)
	if err != nil {
		log.Fatalf("ERROR: Failed adding event to Firestore: %v", err)
		return err
	}

	return c.Status(201).JSON(fiber.Map{"status": "GitHub Event logged successfully."})
}

// Verify GitHub Secret
func verifySignature(secret []byte, data []byte, signature string) bool {
	mac := hmac.New(sha1.New, secret)
	mac.Write(data)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte("sha1="+expectedMAC), []byte(signature))
}
