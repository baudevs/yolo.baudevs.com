package ai

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

// Client handles AI operations
type Client struct {
	client *openai.Client
}

// NewClient creates a new AI client instance
func NewClient(config *Config) *Client {
	apiKey := config.GetAPIKey(config.DefaultProvider)
	return &Client{
		client: openai.NewClient(apiKey),
	}
}

// GetCompletion gets a completion from the AI model
func (c *Client) GetCompletion(prompt string) (string, error) {
	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.7, // More creative for general responses
		},
	)

	if err != nil {
		return "", fmt.Errorf("failed to get completion: %w", err)
	}

	return resp.Choices[0].Message.Content, nil
}
