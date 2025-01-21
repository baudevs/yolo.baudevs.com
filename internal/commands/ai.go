package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"encoding/json"

	"github.com/baudevs/yolo.baudevs.com/internal/ai"
	"github.com/baudevs/yolo.baudevs.com/internal/license"
	"github.com/spf13/cobra"
)

func NewAICommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ai",
		Short: "AI-related commands",
		Long:  `Commands for configuring and using AI features.`,
	}

	cmd.AddCommand(
		newAIConfigureCommand(),
		newAIStatusCommand(),
	)

	return cmd
}

func newAIConfigureCommand() *cobra.Command {
	var apiKey string

	cmd := &cobra.Command{
		Use:   "configure",
		Short: "Configure AI settings",
		Long:  `Configure AI settings like API keys and model preferences.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Initialize config
			config := &ai.Config{
				Model:          "gpt-4",
				DefaultOpenAIKey: apiKey,
				Prompts: map[string]string{
					"commit": "Generate a commit message for these changes:\n\n%s",
					"error":  "Analyze this error:\n\nContext: %s\nError: %v",
					"ask":    "Answer this programming question:\n\n%s",
				},
			}

			// Get config path
			configDir, err := os.UserConfigDir()
			if err != nil {
				return fmt.Errorf("failed to get config directory: %w", err)
			}
			configPath := filepath.Join(configDir, "yolo", "ai_config.json")

			// Create directory if it doesn't exist
			if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
				return fmt.Errorf("failed to create config directory: %w", err)
			}

			// Save config
			data, err := json.MarshalIndent(config, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal config: %w", err)
			}
			if err := os.WriteFile(configPath, data, 0644); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			fmt.Println("✅ AI configuration saved successfully!")
			return nil
		},
	}

	cmd.Flags().StringVarP(&apiKey, "api-key", "k", "", "OpenAI API key")

	return cmd
}

func newAIStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Check AI status",
		Long:  `Check the status of AI features and configuration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Initialize license manager
			manager, err := license.NewManager(license.Config{
				StripeSecretKey:  os.Getenv("STRIPE_SECRET_KEY"),
				DefaultOpenAIKey: os.Getenv("OPENAI_API_KEY"),
			})
			if err != nil {
				return fmt.Errorf("failed to initialize license manager: %w", err)
			}

			// Get license
			lic := manager.GetLicense()

			fmt.Println("🤖 AI Status")
			fmt.Println("-----------")

			// Check if we have an API key
			config := manager.GetConfig()
			if config.DefaultOpenAIKey != "" {
				fmt.Println("✅ Using custom OpenAI API key")
			} else if lic != nil && lic.IsActive {
				fmt.Println("✅ Using YOLO license")
				if lic.PlanType == license.PlanUnlimited {
					fmt.Println("Credits: ♾️  Unlimited")
				} else {
					fmt.Printf("Credits: %d remaining\n", lic.CreditsLeft)
				}
			} else {
				fmt.Println("❌ No active license or API key")
				fmt.Println("\nTo get started, either:")
				fmt.Println("1. Configure your own OpenAI API key:")
				fmt.Println("   yolo ai configure --api-key YOUR_API_KEY")
				fmt.Println("\n2. Purchase YOLO credits:")
				fmt.Println("   yolo license activate")
			}

			return nil
		},
	}
}
