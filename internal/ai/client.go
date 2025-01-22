package ai

import (
	"context"
	"fmt"

	"github.com/baudevs/yolo.baudevs.com/internal/config"
	"github.com/baudevs/yolo.baudevs.com/internal/license"
	"github.com/sashabaranov/go-openai"
)

// Client represents an AI client
type Client struct {
	client         *openai.Client
	licenseManager *license.Manager
}

// NewClient creates a new AI client
func NewClient(cfg *config.ClientConfig, licenseManager *license.Manager) (*Client, error) {
	// Get API key from config or license
	apiKey := cfg.OpenAIKey
	if apiKey == "" {
		// Try to get key from license
		lic := licenseManager.GetLicense()
		if lic == nil || !lic.IsActive {
			return nil, fmt.Errorf("no API key or active license found")
		}
		apiKey = lic.APIKey
	}

	// Create OpenAI client
	client := openai.NewClient(apiKey)

	return &Client{
		client:         client,
		licenseManager: licenseManager,
	}, nil
}

// Ask asks a question to the AI
func (c *Client) Ask(ctx context.Context, query string) (string, error) {
	// Check if we have enough credits
	if err := c.checkCredits(); err != nil {
		return "", err
	}

	// Create chat completion
	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: query,
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to create chat completion: %w", err)
	}

	// Deduct credits if using license
	if err := c.deductCredits(); err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

// GenerateCommitMessage generates a commit message based on the diff
func (c *Client) GenerateCommitMessage(ctx context.Context, diff string) (string, error) {
	// Check if we have enough credits
	if err := c.checkCredits(); err != nil {
		return "", err
	}

	// Create prompt
	prompt := fmt.Sprintf("Generate a concise and descriptive commit message for the following changes:\n\n%s", diff)

	// Create chat completion
	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are a helpful assistant that generates clear and descriptive git commit messages.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to create chat completion: %w", err)
	}

	// Deduct credits if using license
	if err := c.deductCredits(); err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

// checkCredits checks if we have enough credits
func (c *Client) checkCredits() error {
	if c.licenseManager == nil || c.licenseManager.GetLicense() == nil {
		return nil // Using custom API key
	}

	lic := c.licenseManager.GetLicense()
	if !lic.IsActive {
		return fmt.Errorf("no active license")
	}

	if lic.PlanType != "unlimited" && lic.Credits <= 0 {
		return fmt.Errorf("no credits remaining")
	}

	return nil
}

// deductCredits deducts credits for an API call
func (c *Client) deductCredits() error {
	if c.licenseManager == nil || c.licenseManager.GetLicense() == nil {
		return nil // Using custom API key
	}

	return c.licenseManager.DeductCredits(1)
}
