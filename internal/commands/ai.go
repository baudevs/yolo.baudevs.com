package commands

import (
	"fmt"

	"github.com/baudevs/yolo.baudevs.com/internal/ai"
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
	cmd := &cobra.Command{
		Use:   "configure",
		Short: "Configure AI provider settings",
		Long:  "Configure settings for AI providers like OpenAI",
		Run: func(cmd *cobra.Command, args []string) {
			provider, _ := cmd.Flags().GetString("provider")
			apiKey, _ := cmd.Flags().GetString("key")
			model, _ := cmd.Flags().GetString("model")

			if provider == "" || apiKey == "" {
				fmt.Println("Please provide both provider and API key")
				return
			}

			// Load existing config or create new one
			config, err := ai.LoadConfig()
			if err != nil {
				fmt.Printf("Failed to load config: %v\n", err)
				return
			}

			// Update provider configuration
			config.Providers[provider] = ai.Provider{
				Name:    provider,
				APIKey:  apiKey,
				Model:   model,
				Enabled: true,
			}

			// Set as default if it's the only provider
			if len(config.Providers) == 1 {
				config.DefaultProvider = provider
			}

			// Save configuration
			if err := ai.SaveConfig(config); err != nil {
				fmt.Printf("Failed to save configuration: %v\n", err)
				return
			}

			fmt.Printf("Successfully configured %s as AI provider\n", provider)
		},
	}

	cmd.Flags().StringP("provider", "p", "", "AI provider (e.g., openai)")
	cmd.Flags().StringP("key", "k", "", "API key for the provider")
	cmd.Flags().StringP("model", "m", "gpt-3.5-turbo", "Model to use (e.g., gpt-3.5-turbo)")
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

			for _, p := range config.Providers {
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
			msg, _, err := commitAI.GenerateCommitMessage("test: add sample file for testing")
			if err != nil {
				return fmt.Errorf("test failed: %w", err)
			}

			fmt.Println("\n✅ Test successful!")
			fmt.Printf("\nSample commit message:\n%s\n", msg)

			return nil
		},
	}
}
