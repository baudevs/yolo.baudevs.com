package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type AIProvider struct {
	Name     string
	APIKey   string
	BaseURL  string
	Model    string
	Prompts  map[string]string
}

type CommitAI struct {
	Provider AIProvider
	Context  string
	Template string
}

func NewCommitAI() (*CommitAI, error) {
	provider, err := detectProvider()
	if err != nil {
		return nil, err
	}

	return &CommitAI{
		Provider: provider,
		Template: `Analyze the following git changes and generate a conventional commit message.
The message should follow this format:
type(scope): description

Types: feat, fix, docs, style, refactor, test, chore
Scope: optional, describes section of codebase
Description: present-tense summary

Changes:
%s`,
	}, nil
}

func detectProvider() (AIProvider, error) {
	// Try OpenAI
	if key := os.Getenv("OPENAI_API_KEY"); key != "" {
		return AIProvider{
			Name:    "openai",
			APIKey:  key,
			BaseURL: "https://api.openai.com/v1/chat/completions",
			Model:   "gpt-4",
		}, nil
	}

	// Try Anthropic
	if key := os.Getenv("ANTHROPIC_API_KEY"); key != "" {
		return AIProvider{
			Name:    "anthropic",
			APIKey:  key,
			BaseURL: "https://api.anthropic.com/v1/messages",
			Model:   "claude-3-opus-20240229",
		}, nil
	}

	// Try Mistral
	if key := os.Getenv("MISTRAL_API_KEY"); key != "" {
		return AIProvider{
			Name:    "mistral",
			APIKey:  key,
			BaseURL: "https://api.mistral.ai/v1/chat/completions",
			Model:   "mistral-large-latest",
		}, nil
	}

	return AIProvider{}, fmt.Errorf("no AI provider configured")
}

func (c *CommitAI) GenerateCommitMessage(changes string) (string, error) {
	prompt := fmt.Sprintf(c.Template, changes)

	var message string
	var err error

	switch c.Provider.Name {
	case "openai":
		message, err = c.callOpenAI(prompt)
	case "anthropic":
		message, err = c.callAnthropic(prompt)
	case "mistral":
		message, err = c.callMistral(prompt)
	default:
		return "", fmt.Errorf("unsupported AI provider: %s", c.Provider.Name)
	}

	if err != nil {
		return "", err
	}

	// Clean up the message
	message = strings.TrimSpace(message)
	if !strings.Contains(message, ":") {
		return "", fmt.Errorf("invalid commit message format")
	}

	return message, nil
}

func (c *CommitAI) callOpenAI(prompt string) (string, error) {
	payload := map[string]interface{}{
		"model": c.Provider.Model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are a helpful assistant that generates conventional commit messages.",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	return c.makeAPICall(payload)
}

func (c *CommitAI) callAnthropic(prompt string) (string, error) {
	payload := map[string]interface{}{
		"model": c.Provider.Model,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	return c.makeAPICall(payload)
}

func (c *CommitAI) callMistral(prompt string) (string, error) {
	payload := map[string]interface{}{
		"model": c.Provider.Model,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	return c.makeAPICall(payload)
}

func (c *CommitAI) makeAPICall(payload interface{}) (string, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", c.Provider.BaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Provider.APIKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	// Extract message based on provider
	switch c.Provider.Name {
	case "openai":
		choices := result["choices"].([]interface{})
		message := choices[0].(map[string]interface{})["message"].(map[string]interface{})
		return message["content"].(string), nil
	case "anthropic":
		content := result["content"].([]interface{})
		text := content[0].(map[string]interface{})["text"].(string)
		return text, nil
	case "mistral":
		choices := result["choices"].([]interface{})
		message := choices[0].(map[string]interface{})["message"].(map[string]interface{})
		return message["content"].(string), nil
	default:
		return "", fmt.Errorf("unsupported provider response format")
	}
} 