package ai

import (
	"context"
	"fmt"

	"github.com/baudevs/yolo.baudevs.com/internal/license"
	"github.com/baudevs/yolo.baudevs.com/internal/messages"
	openai "github.com/sashabaranov/go-openai"
)

// Client represents an AI client
type Client struct {
	openaiClient *openai.Client
	config       *Config
	license      *license.Manager
	personality  messages.PersonalityLevel
}

// NewClient creates a new AI client
func NewClient(config *Config, licenseManager *license.Manager) (*Client, error) {
	// Get API key in order of priority:
	// 1. User's API key from config
	// 2. Active license API key
	// 3. Default OpenAI key from license manager
	apiKey := config.DefaultOpenAIKey

	if apiKey == "" {
		// Try to get from license
		lic := licenseManager.GetLicense()
		if lic != nil && lic.IsActive {
			apiKey = lic.APIKey
		}
	}

	if apiKey == "" {
		// Finally try license manager config
		apiKey = licenseManager.GetConfig().DefaultOpenAIKey
	}

	if apiKey == "" {
		return nil, fmt.Errorf("no OpenAI API key found. Please run 'yolo license activate' or 'yolo ai configure'")
	}

	return &Client{
		openaiClient: openai.NewClient(apiKey),
		config:       config,
		license:      licenseManager,
		personality:  messages.GetPersonality(),
	}, nil
}

// HasValidLicense checks if there's a valid license or API key
func (c *Client) HasValidLicense() bool {
	// Check if user has their own API key
	if c.config.DefaultOpenAIKey != "" {
		return true
	}

	// Check if user has an active license with credits
	lic := c.license.GetLicense()
	if lic != nil && lic.IsActive {
		if lic.PlanType == license.PlanUnlimited {
			return true
		}
		return lic.CreditsLeft > 0
	}

	return false
}

// checkCredits verifies if the user has enough credits and deducts them
func (c *Client) checkCredits(ctx context.Context) error {
	// Skip credit check if user has their own API key
	if c.config.DefaultOpenAIKey != "" {
		return nil
	}

	// Record AI request (this will check credits and update them)
	if err := c.license.RecordAIRequest(); err != nil {
		return fmt.Errorf("failed to process AI request: %w", err)
	}

	return nil
}

// GenerateCommitMessage generates a commit message using AI
func (c *Client) GenerateCommitMessage(ctx context.Context, changes string) (string, error) {
	if err := c.checkCredits(ctx); err != nil {
		return "", err
	}

	// Get commit prompt based on personality
	var prompt string
	switch c.personality {
	case messages.NerdyClean:
		prompt = fmt.Sprintf(c.config.Prompts["commit"], changes)
	case messages.MildlyRude:
		prompt = fmt.Sprintf("Yo dawg, check out these changes and give me a commit message:\n\n%s\n\nKeep it real but professional-ish.", changes)
	case messages.UnhingedFunny:
		prompt = fmt.Sprintf("YOLO! Time to commit some code! Check this out:\n\n%s\n\nGive me a wild commit message that'll make the team laugh!", changes)
	default:
		prompt = fmt.Sprintf(c.config.Prompts["commit"], changes)
	}

	return c.GenerateResponse(ctx, prompt)
}

// AnalyzeError analyzes an error using AI
func (c *Client) AnalyzeError(ctx context.Context, err error, context string) (string, error) {
	if err := c.checkCredits(ctx); err != nil {
		return "", err
	}

	// Get error prompt based on personality
	var prompt string
	switch c.personality {
	case messages.NerdyClean:
		prompt = fmt.Sprintf(c.config.Prompts["error"], context, err)
	case messages.MildlyRude:
		prompt = fmt.Sprintf("Bruh, we got an error:\n\nContext: %s\nError: %v\n\nWhat's wrong and how do we fix it? Keep it real.", context, err)
	case messages.UnhingedFunny:
		prompt = fmt.Sprintf("CODE RED! ERROR ALERT! 🚨\n\nContext: %s\nError: %v\n\nHelp me fix this disaster! Make it funny!", context, err)
	default:
		prompt = fmt.Sprintf(c.config.Prompts["error"], context, err)
	}

	return c.GenerateResponse(ctx, prompt)
}

// Ask asks a programming question
func (c *Client) Ask(ctx context.Context, question string) (string, error) {
	if err := c.checkCredits(ctx); err != nil {
		return "", err
	}

	// Get ask prompt based on personality
	var prompt string
	switch c.personality {
	case messages.NerdyClean:
		prompt = fmt.Sprintf(c.config.Prompts["ask"], question)
	case messages.MildlyRude:
		prompt = fmt.Sprintf("Yo, here's a question for ya:\n\n%s\n\nBreak it down for me in 3 steps, and don't sugarcoat it.", question)
	case messages.UnhingedFunny:
		prompt = fmt.Sprintf("INCOMING QUESTION ALERT! 🎯\n\n%s\n\nGive me 3 wild steps to solve this, and make it entertaining!", question)
	default:
		prompt = fmt.Sprintf(c.config.Prompts["ask"], question)
	}

	return c.GenerateResponse(ctx, prompt)
}

// GenerateResponse generates a response using the OpenAI API
func (c *Client) GenerateResponse(ctx context.Context, prompt string) (string, error) {
	resp, err := c.openaiClient.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: c.config.Model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return "", fmt.Errorf("failed to generate response: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response generated")
	}

	return resp.Choices[0].Message.Content, nil
}
