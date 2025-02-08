package cmd

import (
	"fmt"

	"github.com/baudevs/yolo.baudevs.com/internal/config"
	"github.com/spf13/cobra"
)

// NewDevCommand returns a new dev command
func NewDevCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dev",
		Short: "Toggle development mode",
		Long:  "Toggle development mode for debugging and testing",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load config
			clientConfig, err := config.LoadClientConfig()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// Toggle dev mode
			clientConfig.DevMode = !clientConfig.DevMode

			// Save config
			if err := config.SaveClientConfig(clientConfig); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			if clientConfig.DevMode {
				fmt.Println(" Development mode enabled")
			} else {
				fmt.Println(" Development mode disabled")
			}

			return nil
		},
	}

	return cmd
}
