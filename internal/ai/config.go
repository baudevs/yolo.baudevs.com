package ai

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"gopkg.in/yaml.v3"
)

// Provider represents an AI provider configuration
type Provider struct {
	Name     string            `yaml:"name"`
	APIKey   string            `yaml:"api_key,omitempty"`
	BaseURL  string            `yaml:"base_url,omitempty"`
	Model    string            `yaml:"model"`
	Prompts  map[string]string `yaml:"prompts,omitempty"`
	Enabled  bool             `yaml:"enabled"`
}

// Config represents the AI configuration
type Config struct {
	DefaultProvider string              `yaml:"default_provider"`
	Providers      map[string]Provider `yaml:"providers"`
}

// LoadConfig loads the AI configuration from the config file
func LoadConfig() (*Config, error) {
	// Get config directory
	configDir, err := getConfigDir()
	if err != nil {
		return nil, err
	}

	// Read config file
	configFile := filepath.Join(configDir, "ai_config.yml")
	data, err := os.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			// Create default config
			return createDefaultConfig(configDir)
		}
		return nil, err
	}

	// Parse config
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}

// SaveConfig saves the AI configuration to the config file
func SaveConfig(config *Config) error {
	// Get config directory
	configDir, err := getConfigDir()
	if err != nil {
		return err
	}

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal config
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write config file
	configFile := filepath.Join(configDir, "ai_config.yml")
	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetAPIKey gets the API key for the specified provider
func (c *Config) GetAPIKey(provider string) string {
	// Check environment variable first
	envKey := fmt.Sprintf("%s_API_KEY", strings.ToUpper(provider))
	if key := os.Getenv(envKey); key != "" {
		return key
	}

	// Fall back to config file
	if p, ok := c.Providers[provider]; ok {
		return p.APIKey
	}

	return ""
}

// createDefaultConfig creates a default AI configuration
func createDefaultConfig(configDir string) (*Config, error) {
	config := &Config{
		DefaultProvider: "openai",
		Providers: map[string]Provider{
			"openai": {
				Name:    "OpenAI",
				Model:   "gpt-3.5-turbo",
				Enabled: true,
				Prompts: map[string]string{
					"commit": defaultCommitPrompt,
					"error":  defaultErrorPrompt,
				},
			},
			"anthropic": {
				Name:    "Anthropic Claude",
				Model:   "claude-2",
				Enabled: false,
			},
			"mistral": {
				Name:    "Mistral AI",
				Model:   "mistral-medium",
				Enabled: false,
			},
		},
	}

	// Save default config
	if err := SaveConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

// getConfigDir gets the YOLO config directory
func getConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return filepath.Join(homeDir, ".yolo"), nil
}

// Default prompts
const (
	defaultCommitPrompt = `Analyze the following Git changes and generate a conventional commit message in JSON format.

Follow these rules:
1. Use semantic commit types: feat, fix, docs, style, refactor, perf, test, build, ci, chore
2. Keep the subject line clear and concise (max 72 chars)
3. Use present tense ("add" not "added")
4. Don't end the subject line with a period
5. Add relevant details in the body
6. Mark breaking changes with breaking=true
7. Include issue references if found in the changes
8. Add co-authors if multiple contributors are detected

The response must be a valid JSON object with this structure:
{
  "type": "feat|fix|docs|style|refactor|perf|test|build|ci|chore",
  "scope": "optional area affected",
  "subject": "concise description",
  "body": "optional detailed explanation",
  "breaking": false,
  "issue_refs": ["optional array of issue references"],
  "co_authors": ["optional array of co-authors"]
}`

	defaultErrorPrompt = `Analyze the following error in the context of a Git operation and provide a helpful explanation and solutions.
Context: %s
Error: %v

Respond with a JSON object containing:
{
  "problem": "brief description of the issue",
  "explanation": "user-friendly explanation of what went wrong",
  "solutions": ["array of step-by-step solutions"]
}`
)
