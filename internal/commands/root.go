package commands

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "yolo",
	Short: "YOLO CLI - Your Optimal Life Organizer",
	Long: `YOLO CLI is a powerful tool that helps you organize your development workflow
with a fun and personality-driven approach.`,
}

// Execute executes the root command
func Execute() error {
	// Add all commands
	rootCmd.AddCommand(InitMessagePromptsCmd())

	return rootCmd.Execute()
}
