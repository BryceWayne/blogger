package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"cloud.google.com/go/firestore"
)

type OpenAI struct {
	APIKey string
}

type OpenAIEmbeddingResponse struct {
	Object string `json:"object"`
	Data   []struct {
		Object    string    `json:"object"`
		Embedding []float64 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
	Model string `json:"model"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

type ChatCompletionRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// GenerateInstructionContent builds a concise doc template.
func GenerateInstructionContent(content, ext string) string {
	const tmpl = `# %s Code Doc\n` +
		"\\`\\`\\`%s\\`\\`\\`\n" +
		"1. Intro\n" +
		"2. Pre-reqs\n" +
		"3. Line-by-Line\n" +
		"4. Techniques if fn, Inputs if struct\n" +
		"5. Examples\n" +
		"6. Summary & TODOs\n" +
		"Return: One paragraph summary of 1-6"

	return fmt.Sprintf(tmpl, ext, content)
}

// GenerateSystemMessage generates the system message for the prompt.
func GenerateSystemMessage() Message {
	return Message{
		Role:    "system",
		Content: "You are a helpful and concise assistant.",
	}
}

// GenerateUserMessage generates the user message for the prompt.
func GenerateUserMessage(instructionContent string) Message {
	return Message{
		Role:    "user",
		Content: instructionContent,
	}
}

// GenerateBlogPrompt creates the blog prompt.
func GenerateBlogPrompt(fileContent string, fileExtension string) []Message {
	instructionContent := GenerateInstructionContent(fileContent, fileExtension)
	systemMessage := GenerateSystemMessage()
	userMessage := GenerateUserMessage(instructionContent)
	return []Message{systemMessage, userMessage}
}

func NewOpenAI(apiKey string) *OpenAI {
	return &OpenAI{APIKey: apiKey}
}

func (o *OpenAI) GenerateBlogPost(file, localFile string, client *firestore.Client) (string, error) {
	// Get file contents

	content, err := ioutil.ReadFile(localFile)
	if err != nil {
		log.Printf("ERROR: Failed to read file %s: %v", localFile, err)
		return "", err
	}

	// Parse file extension
	fileExtension := filepath.Ext(file)
	log.Printf("DEBUG: File extension: %s", fileExtension)

	messages := GenerateBlogPrompt(string(content), fileExtension)

	blogPost, err := o.CreateGPTPrompt(messages)
	if err != nil {
		log.Printf("ERROR: Failed to generate blog post: %v", err)
		return "", err
	}

	return blogPost, nil
}

func (o *OpenAI) UpdateBlogPosts(files []string) {
	for _, file := range files {
		log.Printf("INFO: Updating blog post for file: %s\n", file)
		// Get file contents
		// Update blog post
		// Update blog post in Firestore
	}
}

func (o *OpenAI) RemoveBlogPosts(files []string) {
	for _, file := range files {
		log.Printf("INFO: Removing blog post for file: %s\n", file)
		// Remove blog post from Firestore
	}
}

func (o *OpenAI) CreateGPTPrompt(messages []Message) (string, error) {
	url := "https://api.openai.com/v1/chat/completions"

	reqBody := ChatCompletionRequest{
		Model:    "gpt-4",
		Messages: messages,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("ERROR: Error marshaling request body: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBytes))
	if err != nil {
		return "", fmt.Errorf("ERROR: Error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	// log.Println("API Key:", o.APIKey)
	req.Header.Set("Authorization", "Bearer "+o.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ERROR: Error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ERROR: Error reading response: %v", err)
	}
	log.Println("DEBUG: HTTP Response Status:", resp.Status)
	// log.Println("Response Body:", string(body))

	var response ChatCompletionResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("ERROR: Error unmarshaling response: %v", err)
	}

	if len(response.Choices) > 0 {
		// Your code that uses response.Choices[0]
		return response.Choices[0].Message.Content, nil
	} else {
		// Handle the error or return
		return "", fmt.Errorf("ERROR: no choices returned")
	}

}
