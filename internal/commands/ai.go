package commands

import (
	"fmt"

	"github.com/baudevs/yolo.baudevs.com/internal/config"
	"github.com/baudevs/yolo.baudevs.com/internal/license"
	"github.com/baudevs/yolo.baudevs.com/internal/messages"
	"github.com/spf13/cobra"
)

// NewAICommand returns a new AI command
func NewAICommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ai",
		Short: "AI-related commands",
		Long:  "Configure and manage AI settings",
	}

	cmd.AddCommand(
		newAIConfigCommand(),
		newAIStatusCommand(),
	)

	return cmd
}

func newAIConfigCommand() *cobra.Command {
	var openAIKey string
	var personality string

	cmd := &cobra.Command{
		Use:   "config",
		Short: "Configure AI settings",
		Long:  "Configure AI settings like API key and personality",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load config
			clientConfig, err := config.LoadClientConfig()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// Update OpenAI key if provided
			if openAIKey != "" {
				clientConfig.OpenAIKey = openAIKey
			}

			// Update personality if provided
			if personality != "" {
				level := messages.GetPersonalityFromString(personality)
				if level == messages.Unknown {
					return fmt.Errorf("invalid personality: %s", personality)
				}
				clientConfig.PersonalityType = personality
				messages.SetPersonality(level)
			}

			// Save config
			if err := config.SaveClientConfig(clientConfig); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			fmt.Println("✅ AI settings updated successfully!")
			return nil
		},
	}

	cmd.Flags().StringVarP(&openAIKey, "key", "k", "", "OpenAI API key")
	cmd.Flags().StringVarP(&personality, "personality", "p", "", "AI personality (nerdy, rude, unhinged)")

	return cmd
}

func newAIStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Check AI status",
		Long:  "Check AI configuration status and settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load config
			clientConfig, err := config.LoadClientConfig()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// Create license manager
			licenseManager, err := license.NewManager()
			if err != nil {
				return fmt.Errorf("failed to create license manager: %w", err)
			}

			// Get license
			lic := licenseManager.GetLicense()

			fmt.Println("🤖 AI Status")
			fmt.Println("-----------")

			// Check if we have an API key
			if clientConfig.OpenAIKey != "" {
				fmt.Println("✅ Using custom OpenAI API key")
			} else if lic != nil && lic.IsActive {
				fmt.Println("✅ Using YOLO license")
				if lic.PlanType == "unlimited" {
					fmt.Println("Credits: ♾️  Unlimited")
				} else {
					fmt.Printf("Credits: %d remaining\n", lic.Credits)
				}
			} else {
				fmt.Println("❌ No active license or API key")
				fmt.Println("\nTo get started, either:")
				fmt.Println("1. Configure your own OpenAI API key:")
				fmt.Println("   yolo ai config --key YOUR_API_KEY")
				fmt.Println("\n2. Purchase YOLO credits:")
				fmt.Println("   yolo license activate")
			}

			// Show personality
			fmt.Printf("\nPersonality: %s\n", clientConfig.PersonalityType)

			return nil
		},
	}
}
