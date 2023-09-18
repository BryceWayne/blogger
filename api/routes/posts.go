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

func RecursiveFetch(ctx context.Context, client *firestore.Client, path string, posts *[]models.Post) error {
	iter := client.Collection(path).Documents(ctx)
	for {
		var post models.Post
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("ERROR: Failed to iterate: %v", err)
			return err
		}
		if err := doc.DataTo(&post); err != nil {
			log.Printf("ERROR: Failed to convert document to post: %v", err)
			return err
		}
		log.Printf("DEBUG: Post: %v", post)
		*posts = append(*posts, post)

		// Recurse into nested collections
		collsIter := doc.Ref.Collections(ctx)
		for {
			collRef, err := collsIter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Printf("ERROR: Failed to list collections: %v", err)
				return err
			}
			err = RecursiveFetch(ctx, client, path+"/"+doc.Ref.ID+"/"+collRef.ID, posts)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func GetPosts(c *fiber.Ctx, client *firestore.Client) error {
	ctx := context.Background()
	var posts []models.Post

	err := RecursiveFetch(ctx, client, "posts", &posts)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "Failed to retrieve posts"})
	}

	return c.Status(200).JSON(posts)
}
