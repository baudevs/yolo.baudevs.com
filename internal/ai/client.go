package ai

import (
	"context"
	"fmt"

	"github.com/baudevs/yolo.baudevs.com/internal/config"
	"github.com/baudevs/yolo.baudevs.com/internal/license"
	"github.com/sashabaranov/go-openai"
)

// Client handles communication with the AI service
type Client struct {
	client *openai.Client
}

// NewClient creates a new AI client
func NewClient(cfg *config.Config, licenseManager *license.Manager) (*Client, error) {
	// Try to get API key from config first
	apiKey := cfg.OpenAI.APIKey

	// If not in config, try license manager
	if apiKey == "" {
		var err error
		apiKey, err = licenseManager.GetOpenAIKey()
		if err != nil {
			return nil, fmt.Errorf("failed to get OpenAI key: %w", err)
		}
	}

	// Create OpenAI client
	client := openai.NewClient(apiKey)

	return &Client{
		client: client,
	}, nil
}

// Ask sends a question to the AI and returns the response
func (c *Client) Ask(ctx context.Context, prompt string) (string, error) {
	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4TurboPreview,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return "", fmt.Errorf("failed to create chat completion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	return resp.Choices[0].Message.Content, nil
}

// CreateFunctionCall sends a request to the AI with function definitions and returns the function call response
func (c *Client) CreateFunctionCall(ctx context.Context, messages []openai.ChatCompletionMessage, functions []openai.FunctionDefinition) (string, error) {
	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:     openai.GPT4TurboPreview,
			Messages:  messages,
			Functions: functions,
			FunctionCall: &openai.FunctionCall{
				Name: functions[0].Name,
			},
		},
	)

	if err != nil {
		return "", fmt.Errorf("failed to create function call: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	if resp.Choices[0].Message.FunctionCall == nil {
		return "", fmt.Errorf("no function call in response")
	}

	return resp.Choices[0].Message.FunctionCall.Arguments, nil
}

// GenerateCommitMessage generates a commit message based on the diff
func (c *Client) GenerateCommitMessage(ctx context.Context, diff string) (string, error) {
	prompt := fmt.Sprintf(`Generate a commit message for this diff:

%s

The commit message should:
- Follow the Conventional Commits specification
- Be concise but descriptive
- Focus on the main changes
- Use present tense
- Not exceed 100 characters for the first line

Respond with just the commit message, nothing else.`, diff)

	return c.Ask(ctx, prompt)
}
