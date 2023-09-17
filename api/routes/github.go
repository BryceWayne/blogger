package routes

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"log"
	"os"
	"os/exec"

	"cloud.google.com/go/firestore"
	"github.com/BryceWayne/blogger/models"
	"github.com/BryceWayne/blogger/utils"
	"github.com/gofiber/fiber/v2"
	"gopkg.in/src-d/go-git.v4"
)

func NewGitHubEvent(c *fiber.Ctx, config *utils.Config, client *firestore.Client) error {
	log.Println("INFO: New GitHub Event.")

	payload_ := c.Body()
	signature := c.Get("X-Hub-Signature")

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

	go func() {
		// Initialize the git repo object
		repo, err := git.PlainOpen("/app/api")
		if err != nil {
			log.Printf("ERROR: Failed to open git repo: %v", err)
			return
		}

		// Pull the latest changes
		err = pullChanges(repo, "/app")
		if err != nil {
			log.Printf("ERROR: Failed to pull changes: %v", err)
			return
		}

		log.Println("INFO: Local repo updated successfully.")
	}()

	return c.Status(201).JSON(fiber.Map{"status": "GitHub Event logged successfully."})
}

// Verify GitHub Secret
func verifySignature(secret []byte, data []byte, signature string) bool {
	mac := hmac.New(sha1.New, secret)
	mac.Write(data)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte("sha1="+expectedMAC), []byte(signature))
}

func pullChanges(r *git.Repository, path string) error {
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	err = w.Pull(&git.PullOptions{
		RemoteName: "origin",
		Progress:   os.Stdout,
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return err
	}

	// Navigate to the 'api' directory within your repo
	err = os.Chdir(path + "/api")
	if err != nil {
		log.Printf("ERROR: Failed to change directory: %v", err)
		return err
	}

	// Add changes only from the 'api' directory
	cmd := exec.Command("git", "add", ".")
	err = cmd.Run()
	if err != nil {
		log.Printf("ERROR: Failed to add changes: %v", err)
		return err
	}

	// Commit the changes
	cmd = exec.Command("git", "commit", "-m", "API changes")
	err = cmd.Run()
	if err != nil {
		log.Printf("ERROR: Failed to commit changes: %v", err)
		return err
	}

	return nil
}
