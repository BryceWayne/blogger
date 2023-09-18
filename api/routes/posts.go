package routes

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/BryceWayne/blogger/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

func CreatePost(c *fiber.Ctx, client *firestore.Client) error {
	ctx := context.Background()
	// Create a new post
	post := models.Post{
		ID:        uuid.New().String(),
		Title:     "My new post!",
		Author:    "John Doe",
		Content:   "This is my new post.",
		CreatedAt: time.Now(),
	}

	postsRef := client.Collection("posts")
	docRef := postsRef.Doc(post.ID) // Specify document ID
	_, err := docRef.Set(ctx, post) // Set the document at the specified ID
	if err != nil {
		log.Printf("ERROR: Failed adding post: %v", err)
		return c.Status(500).JSON(fiber.Map{"status": "Failed to create post"})
	}

	log.Printf("INFO: Created post: %s\n", post.ID)

	return c.Status(201).JSON(fiber.Map{"status": "Post created successfully"})
}

func GetPosts(c *fiber.Ctx, client *firestore.Client) error {
	ctx := context.Background()
	var posts []models.Post

	// Fetch posts from Firestore
	iter := client.Collection("posts").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("ERROR: Failed to iterate: %v", err)
			return c.Status(500).JSON(fiber.Map{"status": "Internal Server Error"})
		}
		var post models.Post
		if err := doc.DataTo(&post); err != nil {
			log.Printf("ERROR: Failed to convert to Post model: %v", err)
			return c.Status(500).JSON(fiber.Map{"status": "Internal Server Error"})
		}
		posts = append(posts, post)
	}

	return c.Status(200).JSON(posts)
}
