package routes

import (
	"context"
	"log"
	"sort"

	"cloud.google.com/go/firestore"
	"github.com/BryceWayne/blogger/models"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/iterator"
)

func GetPosts(c *fiber.Ctx, client *firestore.Client) error {
	ctx := context.Background()
	var posts []models.Post

	// Fetch posts from Firestore, including nested documents
	iter := client.Collection("posts").Documents(ctx)
	for {
		var post models.Post

		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			log.Printf("ERROR: Failed to iterate: %v", err)
			return c.Status(500).JSON(fiber.Map{"status": "Failed to retrieve posts"})
		}

		if err := doc.DataTo(&post); err != nil {
			log.Printf("ERROR: Failed to convert document to post: %v", err)
			return c.Status(500).JSON(fiber.Map{"status": "Failed to retrieve posts"})
		}

		posts = append(posts, post)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].File < posts[j].File
	})

	return c.Status(200).JSON(posts)
}
