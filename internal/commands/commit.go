package commands

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/baudevs/yolo.baudevs.com/internal/ai"
	"github.com/baudevs/yolo.baudevs.com/internal/config"
	"github.com/baudevs/yolo.baudevs.com/internal/license"
	"github.com/spf13/cobra"
)

// NewCommitCommand returns a new commit command
func NewCommitCommand() *cobra.Command {
	var message string
	var stageAll bool
	var summarizedDiff bool

	cmd := &cobra.Command{
		Use:   "commit",
		Short: "Create a commit message",
		Long: `Generate an AI-powered commit message based on staged changes.

For large changes, you can use --summarized to send only file names and line counts
to the AI instead of full diffs. This helps when you hit AI token limits.`,
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

			// Stage all changes if requested
			if stageAll {
				if err := stageAllChanges(); err != nil {
					return fmt.Errorf("failed to stage changes: %w", err)
				}
			}

			// Get the changes info
			var diff string
			var diffErr error
			if summarizedDiff {
				diff, diffErr = getSummarizedDiff()
			} else {
				diff, diffErr = getFullDiff()

				// If the diff is too large, warn the user and suggest using summarized mode
				if len(diff) > 12000 { // Conservative limit for GPT-4
					fmt.Println("\n⚠️  Warning: Large diff detected!")
					fmt.Println("The changes are quite extensive, which might affect AI analysis quality.")
					fmt.Println("Consider using --summarized flag for better results.")
					fmt.Println("Proceeding with chunked analysis...")
				}
			}

			if diffErr != nil {
				return fmt.Errorf("failed to get diff: %w", diffErr)
			}

			// Generate commit message if not provided
			if message == "" {
				var genErr error
				message, genErr = client.GenerateCommitMessage(cmd.Context(), diff)
				if genErr != nil {
					return fmt.Errorf("failed to generate commit message: %w", genErr)
				}

				if len(diff) > 12000 {
					fmt.Println("\n⚠️  Note: Some changes were too large and had to be summarized.")
					fmt.Println("The commit message may not reflect all changes in detail.")
					fmt.Println("Use --summarized flag next time for better handling of large changes.")
				}
			}

			// Create commit
			if err := createCommit(message); err != nil {
				return fmt.Errorf("failed to create commit: %w", err)
			}

			fmt.Printf("\n✨ Committed with message:\n%s\n", message)
			return nil
		},
	}

	cmd.Flags().StringVarP(&message, "message", "m", "", "Use provided commit message instead of generating one")
	cmd.Flags().BoolVarP(&stageAll, "all", "a", false, "Stage all changes")
	cmd.Flags().BoolVarP(&summarizedDiff, "summarized", "s", false, "Send only file names and line counts to AI")

	return cmd
}

func stageAllChanges() error {
	cmd := exec.Command("git", "add", "-A")
	return cmd.Run()
}

// getFullDiff returns the complete diff of staged changes
func getFullDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// getSummarizedDiff returns a summary of changes (files and line counts)
func getSummarizedDiff() (string, error) {
	// Get list of changed files
	cmd := exec.Command("git", "diff", "--cached", "--numstat")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	var summary strings.Builder
	summary.WriteString("Changed files summary:\n\n")

	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) >= 3 {
			additions := parts[0]
			deletions := parts[1]
			filename := parts[2]

			// Get file type/extension
			var fileType string
			if idx := strings.LastIndex(filename, "."); idx != -1 {
				fileType = filename[idx+1:]
			} else {
				fileType = "no extension"
			}

			summary.WriteString(fmt.Sprintf("- %s (%s)\n", filename, fileType))
			summary.WriteString(fmt.Sprintf("  +%s lines, -%s lines\n", additions, deletions))
		}
	}

	// Add a summary of the types of files changed
	cmd = exec.Command("git", "diff", "--cached", "--name-only")
	output, err = cmd.Output()
	if err != nil {
		return "", err
	}

	fileTypes := make(map[string]int)
	scanner = bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		filename := scanner.Text()
		ext := "no extension"
		if idx := strings.LastIndex(filename, "."); idx != -1 {
			ext = filename[idx+1:]
		}
		fileTypes[ext]++
	}

	summary.WriteString("\nFile types changed:\n")
	for ext, count := range fileTypes {
		summary.WriteString(fmt.Sprintf("- %s: %d files\n", ext, count))
	}

	return summary.String(), nil
}

func createCommit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	return cmd.Run()
}
