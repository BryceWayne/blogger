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

	log.Println("DEBUG: Received valid payload")

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

	go handleCommits(config, payload.Commits)

	return c.Status(201).JSON(fiber.Map{"status": "GitHub Event logged successfully."})
}

func handleCommits(config *utils.Config, commits []models.Commit) {
	files := map[string][]string{
		"Added":    []string{},
		"Modified": []string{},
		"Removed":  []string{},
	}

	for _, commit := range commits {
		files["Added"] = append(files["Added"], commit.Added...)
		files["Modified"] = append(files["Modified"], commit.Modified...)
		files["Removed"] = append(files["Removed"], commit.Removed...)
	}

	logCommitInfo(files)
	pullChanges()

	ai := models.NewOpenAI(config.OpenAIKey)

	go ai.CreateBlogPosts(files["Added"])
	go ai.UpdateBlogPosts(files["Modified"])
	go ai.RemoveBlogPosts(files["Removed"])
}

func logCommitInfo(files map[string][]string) {
	if len(files["Added"]) > 0 {
		log.Println("DEBUG: Added files:", files["Added"])
	}
	if len(files["Modified"]) > 0 {
		log.Println("DEBUG: Modified files:", files["Modified"])
	}
	if len(files["Removed"]) > 0 {
		log.Println("DEBUG: Removed files:", files["Removed"])
	}
}

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
