package routes

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/BryceWayne/blogger/models"
	"github.com/gofiber/fiber/v2"
)

func NewGitHubEvent(c *fiber.Ctx, client *firestore.Client) error {
	log.Println("INFO: New GitHub Event.")
	ctx := context.Background()
	var payload models.WebhookPayload
	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	log.Printf("DEBUG: Payload %+v\n", payload)

	_, _, err := client.Collection("githubEvents").Add(ctx, payload)
	if err != nil {
		log.Fatalf("Failed adding event to Firestore: %v", err)
		return err
	}

	return c.SendString("Received")
}
