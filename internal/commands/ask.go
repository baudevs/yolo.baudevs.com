package commands

import (
	"fmt"

	"github.com/baudevs/yolo.baudevs.com/internal/ai"
	"github.com/baudevs/yolo.baudevs.com/internal/config"
	"github.com/baudevs/yolo.baudevs.com/internal/license"
	"github.com/spf13/cobra"
)

// NewAskCommand returns a new ask command
func NewAskCommand() *cobra.Command {
	var query string

	cmd := &cobra.Command{
		Use:   "ask",
		Short: "Ask the AI a question",
		Long:  "Ask the AI a programming-related question",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load config
			cfg, err := config.LoadConfig()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// Create license manager
			licenseManager, err := license.NewManager()
			if err != nil {
				return fmt.Errorf("failed to create license manager: %w", err)
			}

			// Create AI client
			client, err := ai.NewClient(cfg, licenseManager)
			if err != nil {
				return fmt.Errorf("failed to create AI client: %w", err)
			}

			// Get query from args if not provided as flag
			if query == "" && len(args) > 0 {
				query = args[0]
			}

			// Ask the question
			response, err := client.Ask(cmd.Context(), query)
			if err != nil {
				return fmt.Errorf("failed to get response: %w", err)
			}

			fmt.Printf("\n%s\n", response)
			return nil
		},
	}

	cmd.Flags().StringVarP(&query, "query", "q", "", "Question to ask")

	return cmd
}
