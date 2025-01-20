package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/baudevs/yolo/internal/ai"
	"github.com/spf13/cobra"
)

func AICmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ai",
		Short: "Manage AI providers and settings",
		Long: `Configure and manage AI providers for YOLO.
You can set up different providers (OpenAI, Anthropic, Mistral)
and customize their settings.`,
	}

	cmd.AddCommand(
		aiConfigureCmd(),
		aiListCmd(),
		aiTestCmd(),
	)

	return cmd
}

func aiConfigureCmd() *cobra.Command {
	var provider string
	var apiKey string
	var model string

	cmd := &cobra.Command{
		Use:   "configure",
		Short: "Configure an AI provider",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load config
			config, err := ai.LoadConfig()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// Get provider
			p, ok := config.Providers[provider]
			if !ok {
				return fmt.Errorf("unknown provider: %s", provider)
			}

			// Update provider settings
			if apiKey != "" {
				p.APIKey = apiKey
			}
			if model != "" {
				p.Model = model
			}
			p.Enabled = true

			config.Providers[provider] = p

			// Save config
			if err := ai.SaveConfig(config); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			fmt.Printf("✅ Configured %s provider\n", p.Name)
			return nil
		},
	}

	cmd.Flags().StringVarP(&provider, "provider", "p", "", "AI provider (openai, anthropic, mistral)")
	cmd.Flags().StringVarP(&apiKey, "api-key", "k", "", "API key for the provider")
	cmd.Flags().StringVarP(&model, "model", "m", "", "Model to use (e.g., gpt-3.5-turbo)")

	cmd.MarkFlagRequired("provider")

	return cmd
}

func aiListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List configured AI providers",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load config
			config, err := ai.LoadConfig()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			fmt.Println("🤖 Configured AI Providers:")
			fmt.Printf("\nDefault provider: %s\n\n", config.DefaultProvider)

			for name, p := range config.Providers {
				status := "❌ Disabled"
				if p.Enabled {
					status = "✅ Enabled"
				}

				fmt.Printf("Provider: %s\n", p.Name)
				fmt.Printf("Status: %s\n", status)
				fmt.Printf("Model: %s\n", p.Model)
				if p.BaseURL != "" {
					fmt.Printf("Base URL: %s\n", p.BaseURL)
				}
				if p.APIKey != "" {
					fmt.Printf("API Key: %s***\n", p.APIKey[:4])
				}
				fmt.Println()
			}

			return nil
		},
	}
}

func aiTestCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "test",
		Short: "Test AI provider configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load config
			config, err := ai.LoadConfig()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// Get default provider
			provider := config.Providers[config.DefaultProvider]
			if !provider.Enabled {
				return fmt.Errorf("default provider %s is not enabled", config.DefaultProvider)
			}

			// Test commit message generation
			commitAI, err := ai.NewCommitAI(config.GetAPIKey(config.DefaultProvider))
			if err != nil {
				return fmt.Errorf("failed to initialize AI: %w", err)
			}

			fmt.Println("🧪 Testing AI provider...")
			msg, err := commitAI.GenerateCommitMessage("test: add sample file for testing")
			if err != nil {
				return fmt.Errorf("test failed: %w", err)
			}

			fmt.Println("\n✅ Test successful!")
			fmt.Printf("\nSample commit message:\n%s\n", msg)

			return nil
		},
	}
}
