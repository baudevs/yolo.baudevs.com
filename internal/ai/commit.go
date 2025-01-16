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
		Template: `Hey there! 👋 Could you help me write a nice commit message for these changes?

I want to follow the "conventional commits" style, which looks like this:
type(area): what changed

The type can be:
✨ feat: for new features
🐛 fix: for bug fixes
📚 docs: for documentation
🎨 style: for visual/formatting changes
♻️  refactor: for code improvements
🧪 test: for adding tests
🔧 chore: for maintenance stuff

The area is optional - it's just what part of the project you worked on.
The description should be clear and start with a present-tense verb.

Here are the changes I made:
%s

Could you write a commit message that:
1. Follows this format
2. Is clear and friendly
3. Captures the main changes
4. Starts with the right type

Thanks! 🙌`,
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
			Prompts: map[string]string{
				"system": "You are a friendly and helpful assistant that writes clear, conventional commit messages. You make technical concepts easy to understand and always maintain a positive, encouraging tone.",
			},
		}, nil
	}

	// Try Anthropic
	if key := os.Getenv("ANTHROPIC_API_KEY"); key != "" {
		return AIProvider{
			Name:    "anthropic",
			APIKey:  key,
			BaseURL: "https://api.anthropic.com/v1/messages",
			Model:   "claude-3-opus-20240229",
			Prompts: map[string]string{
				"system": "You are Claude, a friendly AI that helps write clear commit messages. You make technical concepts approachable and always maintain a helpful, positive tone.",
			},
		}, nil
	}

	// Try Mistral
	if key := os.Getenv("MISTRAL_API_KEY"); key != "" {
		return AIProvider{
			Name:    "mistral",
			APIKey:  key,
			BaseURL: "https://api.mistral.ai/v1/chat/completions",
			Model:   "mistral-large-latest",
			Prompts: map[string]string{
				"system": "You are a friendly AI assistant that helps write clear commit messages. You make technical concepts easy to understand and always keep a positive, encouraging tone.",
			},
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