package models

import (
	"bytes"
	"context"
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

func GenerateBlogPrompt(fileContent string, fileExtension string) []Message {
	// Escaping backticks inside the string
	instructionContent := fmt.Sprintf("# Blog Post Creation Based on a Code Snippet\n\n"+
		"## Introduction\n"+
		"Create a detailed and engaging blog post that explains the following code snippet. "+
		"The code is written in %s, and it has a docstring that provides a brief explanation of what the function does.\n\n"+
		"## Code Snippet\n"+
		"\\`\\`\\`%s\n"+
		"%s\n"+
		"\\`\\`\\`\n\n"+
		"## Requirements\n"+
		"1. Start with an engaging introduction that sets the context for the code snippet.\n"+
		"2. Explain any pre-requisites or concepts that the reader should understand before diving into the code.\n"+
		"3. Walk through the code snippet, line-by-line, explaining what each line does.\n"+
		"4. If the function uses any special programming techniques, elaborate on them.\n"+
		"5. Provide one or more use-cases or examples demonstrating how this function could be used in a real-world scenario.\n"+
		"6. Conclude with a summary and potential future enhancements or applications.\n"+
		"7. Use subheadings to break up the text and make it easier to read.\n\n"+
		"## Output\n"+
		"The final blog post should be formatted in Markdown and be both informative and engaging.", "Go", fileExtension, fileContent)

	messages := []Message{
		{
			Role:    "system",
			Content: "You are a helpful assistant.",
		},
		{
			Role:    "user",
			Content: instructionContent,
		},
	}

	return messages
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

func NewOpenAI(apiKey string) *OpenAI {
	return &OpenAI{APIKey: apiKey}
}

func (o *OpenAI) FetchEmbedding(inputText string) (OpenAIEmbeddingResponse, error) {
	url := "https://api.openai.com/v1/embeddings"

	payload := map[string]interface{}{
		"input": inputText,
		"model": "text-embedding-ada-002",
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return OpenAIEmbeddingResponse{}, fmt.Errorf("Error marshaling payload: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return OpenAIEmbeddingResponse{}, fmt.Errorf("Error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return OpenAIEmbeddingResponse{}, fmt.Errorf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return OpenAIEmbeddingResponse{}, fmt.Errorf("Error reading response: %v", err)
	}

	var response OpenAIEmbeddingResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return OpenAIEmbeddingResponse{}, fmt.Errorf("Error unmarshaling response: %v", err)
	}

	return response, nil
}

func (o *OpenAI) GenerateBlogPost(file string, client *firestore.Client) (string, error) {
	// Get file contents
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Printf("ERROR: Failed to read file %s: %v", file, err)
		return "", err
	}

	// Parse file extension
	fileExtension := filepath.Ext(file)
	log.Printf("DEBUG: File extension: %s", fileExtension)

	// Generate blog post using GPT (replace with your actual GPT call)
	messages := GenerateBlogPrompt(string(content), fileExtension)

	blogPost, err := o.CreateGPTPrompt(messages)
	if err != nil {
		log.Printf("ERROR: Failed to generate blog post: %v", err)
		return "", err
	}

	post := struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		File    string `json:"file"`
	}{
		Title:   "Generated title",
		Content: blogPost,
		File:    file,
	}

	docRef := client.Collection("posts").Doc(post.File)
	_, err = docRef.Set(context.Background(), post)
	if err != nil {
		log.Printf("ERROR: Failed to add blog post to Firestore: %v", err)
	} else {
		log.Println("INFO: Successfully added blog post to Firestore.")
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
		Model:    "gpt-3.5-turbo",
		Messages: messages,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("Error marshaling request body: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBytes))
	if err != nil {
		return "", fmt.Errorf("Error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	// log.Println("API Key:", o.APIKey)
	req.Header.Set("Authorization", "Bearer "+o.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading response: %v", err)
	}
	fmt.Println("HTTP Status:", resp.Status)
	fmt.Println("Response Body:", string(body))

	var response ChatCompletionResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("Error unmarshaling response: %v", err)
	}

	if len(response.Choices) > 0 {
		// Your code that uses response.Choices[0]
		return response.Choices[0].Message.Content, nil
	} else {
		// Handle the error or return
		return "", fmt.Errorf("Error: no choices returned")
	}

}
