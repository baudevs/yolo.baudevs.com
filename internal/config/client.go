package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ClientConfig holds client-side configuration
type ClientConfig struct {
	APIEndpoint     string            `json:"api_endpoint"`
	OpenAIKey       string            `json:"openai_key"`
	DevMode         bool              `json:"dev_mode"`
	Prompts         map[string]string `json:"prompts"`
	PersonalityType string            `json:"personality_type"`
}

// Default values
const (
	DefaultAPIEndpoint    = "https://backend.yolodev.dev/api"
	DefaultDevAPIEndpoint = "http://localhost:3000/api"
)

// LoadClientConfig loads the client configuration
func LoadClientConfig() (*ClientConfig, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	// If config doesn't exist, create default
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config := &ClientConfig{
			APIEndpoint:     DefaultAPIEndpoint,
			DevMode:         false,
			Prompts:         make(map[string]string),
			PersonalityType: "nerdy",
		}
		return config, nil
	}

	// Read existing config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config ClientConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Initialize maps if nil
	if config.Prompts == nil {
		config.Prompts = make(map[string]string)
	}

	// Set default personality if empty
	if config.PersonalityType == "" {
		config.PersonalityType = "nerdy"
	}

	// Override API endpoint if in dev mode
	if config.DevMode {
		config.APIEndpoint = DefaultDevAPIEndpoint
	}

	return &config, nil
}

// SaveClientConfig saves the client configuration
func SaveClientConfig(config *ClientConfig) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Save config
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// getConfigPath returns the path to the config file
func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return filepath.Join(home, ".yolo", "config.json"), nil
}
