package commands

import (
	"fmt"
	"os/exec"

	"github.com/baudevs/yolo.baudevs.com/internal/ai"
	"github.com/baudevs/yolo.baudevs.com/internal/config"
	"github.com/baudevs/yolo.baudevs.com/internal/license"
	"github.com/spf13/cobra"
)

// NewCommitCommand returns a new commit command
func NewCommitCommand() *cobra.Command {
	var message string
	var stageAll bool

	cmd := &cobra.Command{
		Use:   "commit",
		Short: "Create a commit message",
		Long:  "Generate an AI-powered commit message based on staged changes",
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

			// Create AI client
			client, err := ai.NewClient(clientConfig, licenseManager)
			if err != nil {
				return fmt.Errorf("failed to create AI client: %w", err)
			}

			// Stage all changes if requested
			if stageAll {
				if err := stageAllChanges(); err != nil {
					return fmt.Errorf("failed to stage changes: %w", err)
				}
			}

			// Get the diff
			diff, err := getDiff()
			if err != nil {
				return fmt.Errorf("failed to get diff: %w", err)
			}

			// Generate commit message if not provided
			if message == "" {
				message, err = client.GenerateCommitMessage(cmd.Context(), diff)
				if err != nil {
					return fmt.Errorf("failed to generate commit message: %w", err)
				}
			}

			// Create commit
			if err := createCommit(message); err != nil {
				return fmt.Errorf("failed to create commit: %w", err)
			}

			fmt.Printf(" Committed with message:\n%s\n", message)
			return nil
		},
	}

	cmd.Flags().StringVarP(&message, "message", "m", "", "Use provided commit message instead of generating one")
	cmd.Flags().BoolVarP(&stageAll, "all", "a", false, "Stage all changes")

	return cmd
}

func stageAllChanges() error {
	cmd := exec.Command("git", "add", "-A")
	return cmd.Run()
}

func getDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func createCommit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	return cmd.Run()
}