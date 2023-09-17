package routes

import (
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/BryceWayne/tooter/models"
	"github.com/BryceWayne/tooter/utils"
	"github.com/gofiber/fiber/v2"
)

// @Summary Generate Content
// @Description Generate content based on user interests
// @ID generate-content
// @Accept  json
// @Produce json
// @Success 200 {object} map[string]interface{} "status: success, content: generated content"
// @Router /generate [get]
func GenerateContent(c *fiber.Ctx, client *firestore.Client) error {
	log.Println("INFO: Generating content.")
	system_prompt := fmt.Sprintf("Be concise. Aim for engagement: likes or replies. Stats: %d likes, %d replies. Goal: More engagement, no AI reveal. Limit: <= 69 chars. Include: variation of original query in response.", likes, replies)
	log.Printf("system_prompt: %s\n", system_prompt)

	var messages []models.Message
	messages = append(messages, models.Message{
		Role:    "system",
		Content: system_prompt,
	})
	messages = append(messages, models.Message{
		Role:    "user",
		Content: utils.GetRandomPrompt(user_prompt),
	})

	// Create GPT-3 prompt based on user interests
	generatedContent, err := openaiClient.CreateGPTPrompt(messages)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create GPT-3 prompt"})
	}

	return c.JSON(fiber.Map{"status": "success", "content": generatedContent})
}
