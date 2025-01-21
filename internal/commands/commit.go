package commands

import (
	"context"
	"fmt"
	"os"
	"strings"
	"path/filepath"
	"encoding/json"

	"github.com/baudevs/yolo.baudevs.com/internal/ai"
	"github.com/baudevs/yolo.baudevs.com/internal/git"
	"github.com/baudevs/yolo.baudevs.com/internal/license"
	"github.com/spf13/cobra"
)

func NewCommitCommand() *cobra.Command {
	var (
		message string
		all     bool
		retry   bool
	)

	commitCmd := &cobra.Command{
		Use:   "commit",
		Short: "Commit changes with AI-generated messages",
		Long: `Commit your changes with AI-generated commit messages.
The AI will analyze your changes and generate a descriptive commit message.

Examples:
  yolo commit              # Generate message for staged changes
  yolo commit -a          # Stage all changes and generate message
  yolo commit -m "feat: add user auth"  # Use custom message`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			// Initialize git client
			gitClient, err := git.NewClient()
			if err != nil {
				return fmt.Errorf("failed to initialize git client: %w", err)
			}

			// Check if we have any changes
			if !gitClient.HasChanges() {
				return fmt.Errorf("no changes to commit")
			}

			// Stage all changes if requested
			if all {
				if err := gitClient.StageAll(); err != nil {
					return fmt.Errorf("failed to stage changes: %w", err)
				}
			}

			// Get staged changes
			diff, err := gitClient.GetStagedDiff()
			if err != nil {
				return fmt.Errorf("failed to get staged changes: %w", err)
			}

			// If no staged changes, error out
			if diff == "" {
				return fmt.Errorf("no staged changes. Use -a to stage all changes")
			}

			// Initialize AI client
			aiClient, err := initAIClient()
			if err != nil {
				return fmt.Errorf("failed to initialize AI client: %w", err)
			}

			// If message is provided, use it directly
			if message != "" {
				if err := gitClient.Commit(message); err != nil {
					// If commit fails, try to analyze error
					analysis, err2 := aiClient.AnalyzeError(ctx, err, "git commit")
					if err2 == nil && analysis != "" {
						return fmt.Errorf("%v\n\nAI Analysis: %s", err, analysis)
					}
					return fmt.Errorf("failed to commit: %w", err)
				}
				return nil
			}

			// Generate commit message
			msg, err := aiClient.GenerateCommitMessage(ctx, diff)
			if err != nil {
				// Try to analyze error
				analysis, err2 := aiClient.AnalyzeError(ctx, err, "generating commit message")
				if err2 == nil && analysis != "" {
					return fmt.Errorf("%v\n\nAI Analysis: %s", err, analysis)
				}
				return fmt.Errorf("failed to generate commit message: %w", err)
			}

			// Commit changes
			if err := gitClient.Commit(msg); err != nil {
				return fmt.Errorf("failed to commit: %w", err)
			}

			fmt.Printf(" Committed with message:\n%s\n", msg)
			return nil
		},
	}

	commitCmd.Flags().StringVarP(&message, "message", "m", "", "Use the given message instead of generating one")
	commitCmd.Flags().BoolVarP(&all, "all", "a", false, "Stage all changes before committing")
	commitCmd.Flags().BoolVarP(&retry, "retry", "r", false, "Retry the last failed commit")

	return commitCmd
}

func initAIClient() (*ai.Client, error) {
	// Initialize license manager
	manager, err := license.NewManager(license.Config{
		StripeSecretKey:  os.Getenv("STRIPE_SECRET_KEY"),
		DefaultOpenAIKey: os.Getenv("OPENAI_API_KEY"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize license manager: %w", err)
	}

	// Load AI config
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config directory: %w", err)
	}
	configPath := filepath.Join(configDir, "yolo", "ai_config.json")

	var config ai.Config
	data, err := os.ReadFile(configPath)
	if err == nil {
		if err := json.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse config: %w", err)
		}
	}

	// Initialize AI client
	aiClient, err := ai.NewClient(&config, manager)
	if err != nil {
		// Try to analyze error
		if strings.Contains(err.Error(), "OpenAI API key") {
			analysis, err2 := aiClient.AnalyzeError(context.Background(), err, "initializing AI client")
			if err2 == nil && analysis != "" {
				return nil, fmt.Errorf("%v\n\nAI Analysis: %s", err, analysis)
			}
		}
		return nil, fmt.Errorf("failed to initialize AI client: %w", err)
	}

	return aiClient, nil
}