package ai

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	configFileName = ".yolo/ai_config.json"
)

// Config represents the AI configuration
type Config struct {
	DefaultOpenAIKey string            `yaml:"default_openai_key" json:"default_openai_key"` // Our API key
	APIKeys          map[string]string  `yaml:"api_keys" json:"api_keys"`
	Model            string            `yaml:"model"`
	Prompts          map[string]string `yaml:"prompts,omitempty"`
}

// LoadConfig loads the AI configuration from disk
func LoadConfig() (*Config, error) {
	// Try to get API key from environment first
	apiKey := os.Getenv("OPENAI_API_KEY")

	// Try to load from YAML config first
	configDir, err := getConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config directory: %w", err)
	}

	var config Config

	// Try to load YAML config
	yamlPath := filepath.Join(configDir, "config.yaml")
	yamlData, err := os.ReadFile(yamlPath)
	if err == nil {
		if err := yaml.Unmarshal(yamlData, &config); err != nil {
			return nil, fmt.Errorf("failed to parse YAML config: %w", err)
		}
	}

	// Try to load JSON config
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	jsonPath := filepath.Join(homeDir, configFileName)
	jsonData, err := os.ReadFile(jsonPath)
	if err == nil {
		if err := json.Unmarshal(jsonData, &config); err != nil {
			return nil, fmt.Errorf("failed to parse JSON config: %w", err)
		}
	}

	// Initialize if empty
	if config.APIKeys == nil {
		config.APIKeys = make(map[string]string)
	}

	// Environment variables take precedence
	if apiKey != "" {
		config.APIKeys["openai"] = apiKey
		config.DefaultOpenAIKey = apiKey
	}

	return &config, nil
}

// SaveConfig saves the configuration to both YAML and JSON formats
func (c *Config) SaveConfig() error {
	// Save YAML config
	configDir, err := getConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config directory: %w", err)
	}

	yamlData, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML config: %w", err)
	}

	yamlPath := filepath.Join(configDir, "config.yaml")
	if err := os.WriteFile(yamlPath, yamlData, 0644); err != nil {
		return fmt.Errorf("failed to write YAML config: %w", err)
	}

	// Save JSON config
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	jsonConfigDir := filepath.Join(homeDir, ".yolo")
	if err := os.MkdirAll(jsonConfigDir, 0755); err != nil {
		return fmt.Errorf("failed to create JSON config directory: %w", err)
	}

	jsonData, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON config: %w", err)
	}

	jsonPath := filepath.Join(homeDir, configFileName)
	if err := os.WriteFile(jsonPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write JSON config: %w", err)
	}

	return nil
}

// GetAPIKey returns the API key for the given provider
func (c *Config) GetAPIKey(provider string) string {
	// Check environment variable first
	envKey := fmt.Sprintf("%s_API_KEY", strings.ToUpper(provider))
	if key := os.Getenv(envKey); key != "" {
		return key
	}

	// Fall back to config file
	if key, ok := c.APIKeys[provider]; ok {
		return key
	}

	// Finally, try default OpenAI key
	if provider == "openai" {
		return c.DefaultOpenAIKey
	}

	return ""
}

// createDefaultConfig creates a default AI configuration
func createDefaultConfig() (*Config, error) {
	config := &Config{
		Model: "gpt-4",
		Prompts: map[string]string{
			"commit": `Analyze these changes and generate a conventional commit message:

%s

The commit message should follow the Conventional Commits specification and include:
1. Type (feat, fix, docs, style, refactor, test, chore)
2. Optional scope in parentheses
3. Description
4. Optional body with breaking changes
5. Optional footer with issue references

Keep it clear and concise.`,
			"error": `Error context: %s
Error message: %v

Please analyze this error and provide:
1. A clear explanation of what went wrong
2. The likely root cause
3. Suggested steps to fix it
4. Any preventive measures for the future

Format the response in a clear, concise way.`,
			"ask": `You are a helpful programming assistant. Please provide a 3-step solution to this question:
%s

Format your response as:
"Here's how to solve that, fellow developer!

1. [First step]
2. [Second step]
3. [Third step]"

Keep each step concise and practical.`,
		},
		APIKeys: make(map[string]string),
	}

	// Create config directory
	configDir, err := getConfigDir()
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	// Save default config
	if err := config.SaveConfig(); err != nil {
		return nil, err
	}

	return config, nil
}

// getConfigDir returns the configuration directory path
func getConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return filepath.Join(homeDir, ".yolo"), nil
}
