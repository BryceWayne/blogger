package routes

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"log"
	"os/exec"

	"cloud.google.com/go/firestore"
	"github.com/BryceWayne/blogger/models"
	"github.com/BryceWayne/blogger/utils"
	"github.com/gofiber/fiber/v2"
)

func NewGitHubEvent(c *fiber.Ctx, config *utils.Config, client *firestore.Client) error {
	log.Println("INFO: New GitHub Event.")

	payload_ := c.Body()
	signature := c.Get("X-Hub-Signature")

	secret := []byte(config.WebhookSecret) // Replace with your GitHub webhook secret

	if !verifySignature(secret, payload_, signature) {
		return c.Status(401).SendString("ERROR: Mismatched signature")
	}

	// log.Println("DEBUG: Received valid payload")

	ctx := context.Background()
	var payload models.WebhookPayload
	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	_, _, err := client.Collection("githubEvents").Add(ctx, payload)
	if err != nil {
		log.Printf("ERROR: Failed adding event to Firestore: %v", err)
		return err
	}
	log.Println("INFO: Added GitHub Event to Firestore.")

	go handleCommits(config, client, payload.Commits)

	return c.Status(201).JSON(fiber.Map{"status": "GitHub Event logged successfully."})
}

func handleCommits(config *utils.Config, client *firestore.Client, commits []models.Commits) {
	// ctx := context.Background()

	fileStatusMap := make(map[string]models.FileStatus)

	for _, commit := range commits {
		for _, addedFile := range commit.Added {
			fileStatus := fileStatusMap[addedFile]
			fileStatus.Added = true
			fileStatus.Modified = false
			fileStatus.Removed = false
			fileStatusMap[addedFile] = fileStatus
			log.Printf("DEBUG: Added file: %s", addedFile)
		}

		for _, modifiedFile := range commit.Modified {
			fileStatus := fileStatusMap[modifiedFile]
			fileStatus.Modified = true
			fileStatus.Added = false
			fileStatus.Removed = false
			fileStatusMap[modifiedFile] = fileStatus
			log.Printf("DEBUG: Modified file: %s", modifiedFile)
		}

		for _, removedFile := range commit.Removed {
			fileStatus := fileStatusMap[removedFile]
			fileStatus.Removed = true
			fileStatus.Added = false
			fileStatus.Modified = false
			fileStatusMap[removedFile] = fileStatus
			log.Printf("DEBUG: Removed file: %s", removedFile)
		}
	}

	// logCommitInfo(files)
	pullChanges()

	ai := models.NewOpenAI(config.OpenAIKey)

	// Now you can iterate through the fileStatusMap and process each file based on its status.
	for filename, status := range fileStatusMap {
		localpath := "/app/blogger/" + filename
		if status.Added {
			// Process added file

			createPostByFileName(localpath, filename, client, ai)

		}
		// if status.Modified {
		// 	// Process modified file
		// 	// Modify the post using the file name (file) and the updated content
		// 	blogPost, err := ai.GenerateBlogPost(file, filepath, client)
		// 	if err != nil {
		// 		log.Printf("Error creating blog post for: %s", file)
		// 		break
		// 	}

		// 	if err := updatePostByFileName(client, file, blogPost); err != nil {
		// 		log.Printf("ERROR: Failed to update post for file %s: %v", file, err)
		// 	} else {
		// 		log.Printf("INFO: Successfully updated post for file %s", file)
		// 	}

		// }
		// if status.Removed {
		// 	// Process removed file
		// 	// Assuming you have a function to remove a post based on the file name
		// 	if err := removePostByFileName(client, file); err != nil {
		// 		log.Printf("ERROR: Failed to remove post for file %s: %v", file, err)
		// 	} else {
		// 		log.Printf("INFO: Successfully removed post for file %s", file)
		// 	}
		// }
	}

}

func createPostByFileName(localpath, filename string, client *firestore.Client, ai *models.OpenAI) {
	// Create a new post using the file name and content
	// ctx := context.Background()

	codeMapInterface := utils.ParseGoFile(localpath)

	// var content string

	// Check if the returned value is a map[string]map[string]string
	codeMap, ok := codeMapInterface.(map[string]map[string]string)
	if !ok {
		log.Println("DEBUG: ParseGoFile didn't return the expected map type")
		return
	}
	log.Printf("DEBUG: Code: %+v\n", codeMap)

	for topLevelKey, topLevelValue := range codeMap {
		log.Printf("DEBUG: Top-Level Key: %s", topLevelKey)
		log.Printf("DEBUG: Top-Level Value: %v", topLevelValue)

		// Now, iterate over the nested map
		for nestedKey, nestedValue := range topLevelValue {
			log.Printf("DEBUG: Nested Key: %s", nestedKey)
			log.Printf("DEBUG: Nested Value: %s", nestedValue)
		}
	}

	// for k, v := range codeMap {
	// 	log.Printf("DEBUG: Key: %s", k)
	// 	log.Printf("DEBUG: Value: %s", v)

	// blogPost, err := ai.GenerateBlogPost(filename, v, client)
	// if err != nil {
	// 	log.Printf("Error creating blog post for code chunk: %s", v)
	// 	break
	// }
	// }

	// post := models.NewPost(filename, "Blog Post", "OpenAI", content)
	// posts := []models.Post{post}

	// docRef := client.Collection("posts").Doc(post.ID)
	// _, err := docRef.Set(ctx, posts)
	// if err != nil {
	// 	log.Printf("ERROR: Failed to add blog post to Firestore: %v", err)
	// 	return
	// }

	log.Printf("INFO: Successfully created post for file %s", filename)
	return
}

// func updatePostByFileName(client *firestore.Client, fileName string, updatedContent string) error {
// 	// Fetch all root toots sorted by LikeCount
// 	ctx := context.Background()

// 	postsRef := client.Collection("posts")

// 	iter := postsRef.Where("File", "==", fileName).Documents(ctx)
// 	defer iter.Stop()

// 	var found bool = false

// 	for {
// 		doc, err := iter.Next()
// 		if err == iterator.Done {
// 			log.Printf("DEBUG: No more posts for file: %s", fileName)
// 			break
// 		}
// 		if err != nil {
// 			log.Printf("ERROR: Failed to fetch post: %v", err)
// 			return fmt.Errorf("Failed to fetch post: %v", err)
// 		}

// 		var post models.Post
// 		if err := doc.DataTo(&post); err != nil {
// 			log.Printf("ERROR: Failed to convert to Post model: %v", err)
// 			return fmt.Errorf("Failed to convert to Post model: %v", err)
// 		}

// 		post.Content = updatedContent
// 		post.UpdatedAt = time.Now()

// 		_, err = doc.Ref.Set(ctx, post)
// 		if err != nil {
// 			log.Printf("ERROR: Failed to update post: %v", err)
// 			return fmt.Errorf("Failed to update post: %v", err)
// 		}

// 		found = true
// 	}

// 	if !found {
// 		err := createPostByFileName(client, fileName, updatedContent)
// 		if err != nil {
// 			log.Printf("ERROR: Failed to create post for file %s: %v", fileName, err)
// 			return fmt.Errorf("Failed to create post for file %s: %v", fileName, err)
// 		}

// 		log.Printf("INFO: Successfully created post for file %s", fileName)
// 	}

// 	return nil
// }

// func removePostByFileName(client *firestore.Client, fileName string) error {
// 	// Fetch all root toots sorted by LikeCount
// 	ctx := context.Background()

// 	postsRef := client.Collection("posts")

// 	iter := postsRef.Where("File", "==", fileName).Documents(ctx)
// 	defer iter.Stop()

// 	for {
// 		doc, err := iter.Next()
// 		if err == iterator.Done {
// 			log.Printf("DEBUG: No more posts for file: %s", fileName)
// 			break
// 		}
// 		if err != nil {
// 			log.Printf("ERROR: Failed to fetch post: %v", err)
// 			return fmt.Errorf("Failed to fetch post: %v", err)
// 		}

// 		_, err = doc.Ref.Delete(ctx)
// 		if err != nil {
// 			log.Printf("ERROR: Failed to delete post: %v", err)
// 			return fmt.Errorf("Failed to delete post: %v", err)
// 		}

// 	}

// 	return nil
// }

// func logCommitInfo(fileStatusMap map[string]models.FileStatus) {
// 	for filePath, fileStatus := range fileStatusMap {
// 		if len(fileStatus.Added) > 0 {
// 			log.Printf("DEBUG: Added files in commit for %s: %v\n", filePath, fileStatus.Added)
// 		}
// 		if len(fileStatus.Modified) > 0 {
// 			log.Printf("DEBUG: Modified files in commit for %s: %v\n", filePath, fileStatus.Modified)
// 		}
// 		if len(fileStatus.Removed) > 0 {
// 			log.Printf("DEBUG: Removed files in commit for %s: %v\n", filePath, fileStatus.Removed)
// 		}
// 	}
// }

func pullChanges() {
	// Shell command to pull changes from GitHub
	cmd := exec.Command("git", "-C", "/app/blogger", "pull", "origin", "master")
	err := cmd.Run()
	if err != nil {
		log.Printf("ERROR: Failed to pull changes from GitHub: %v", err)
	} else {
		log.Println("INFO: Successfully pulled changes from GitHub.")
	}
}

// Verify GitHub Secret
func verifySignature(secret []byte, data []byte, signature string) bool {
	mac := hmac.New(sha1.New, secret)
	mac.Write(data)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte("sha1="+expectedMAC), []byte(signature))
}
