package commands

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// SprintCmd represents the sprint command
func SprintCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sprint",
		Short: "Manage AI-user collaboration sprints",
		Long: `Track and manage the thought process and collaboration between the AI and user during development.

This command helps you document the AI's thought process and planning during development.
It creates and maintains a sprint.current.md file in your project root.`,
	}

	cmd.AddCommand(sprintCurrentCmd())
	return cmd
}

// sprintCurrentCmd represents the sprint current subcommand
func sprintCurrentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "current",
		Short: "Manage the current sprint",
		Long: `Track the current sprint's thought process and collaboration between AI and user.

The 'current' subcommand manages a sprint.current.md file that contains the AI's
thought process, plans, and explanations during development. Use 'init' to start
a new sprint and 'update' to add more information as you progress.`,
	}

	cmd.AddCommand(sprintCurrentInitCmd())
	cmd.AddCommand(sprintCurrentUpdateCmd())
	return cmd
}

// sprintCurrentInitCmd initializes a new sprint.current.md file
func sprintCurrentInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize a new sprint with initial thought process",
		Long: `Create a new sprint.current.md file with the initial thought process from the AI.

There are two ways to provide the content:

1. As a command line argument (for short, single-line text):
   yolo sprint current init "Initial plan: create user authentication"

2. By pasting multi-line text (recommended for longer content):
   yolo sprint current init
   - The command will wait for your input
   - Paste or type your text (it will preserve all formatting)
   - Press Enter for new lines
   - When done, press:
     • Ctrl+D (on Mac/Linux)
     • Ctrl+Z followed by Enter (on Windows)

Example of multi-line input:
   $ yolo sprint current init
   # Implementation Plan
   1. First, we'll create the database schema
   2. Then, add the API endpoints
   
   Example code:
   ` + "```" + `python
   def authenticate():
       pass
   ` + "```" + `
   <press Ctrl+D when done>

The command will create sprint.current.md in your project root with the provided
content and a timestamp. If the file already exists, use 'update' instead.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var content string
			
			// If no args provided, read from stdin
			if len(args) == 0 {
				fmt.Println("📝 Paste or type your text (press Ctrl+D when done):")
				bytes, err := io.ReadAll(os.Stdin)
				if err != nil {
					return fmt.Errorf("failed to read from stdin: %w", err)
				}
				content = string(bytes)
			} else {
				content = strings.Join(args, " ")
			}

			if strings.TrimSpace(content) == "" {
				return fmt.Errorf("no content provided")
			}

			// Get project root
			root, err := getProjectRoot()
			if err != nil {
				return fmt.Errorf("failed to get project root: %w", err)
			}

			// Create sprint file
			sprintFile := filepath.Join(root, "sprint.current.md")
			
			// Check if file already exists
			if _, err := os.Stat(sprintFile); err == nil {
				return fmt.Errorf("sprint.current.md already exists. Use 'update' to add new thoughts")
			}

			// Format content with timestamp
			formattedContent := formatSprintEntry(content)

			// Write to file
			if err := os.WriteFile(sprintFile, []byte(formattedContent), 0644); err != nil {
				return fmt.Errorf("failed to write sprint file: %w", err)
			}

			fmt.Printf("✨ Created new sprint.current.md in %s\n", root)
			return nil
		},
	}
}

// sprintCurrentUpdateCmd appends new thought process to sprint.current.md
func sprintCurrentUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Update current sprint with new thought process",
		Long: `Append new thought process to the sprint.current.md file.

There are two ways to provide the content:

1. As a command line argument (for short, single-line text):
   yolo sprint current update "Added error handling to auth module"

2. By pasting multi-line text (recommended for longer content):
   yolo sprint current update
   - The command will wait for your input
   - Paste or type your text (it will preserve all formatting)
   - Press Enter for new lines
   - When done, press:
     • Ctrl+D (on Mac/Linux)
     • Ctrl+Z followed by Enter (on Windows)

Example of multi-line input:
   $ yolo sprint current update
   ## Progress Update
   - Implemented user authentication
   - Added unit tests
   - Fixed bug in password validation
   
   Next steps:
   1. Add password reset functionality
   2. Implement email verification
   
   Code changes:
   ` + "```" + `python
   def reset_password(user_id: str):
       # TODO: Implement this
       pass
   ` + "```" + `
   <press Ctrl+D when done>

The command will append your update to sprint.current.md with a timestamp and a
clear separator. If the file doesn't exist, use 'init' first.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var content string
			
			// If no args provided, read from stdin
			if len(args) == 0 {
				fmt.Println("📝 Paste or type your text (press Ctrl+D when done):")
				bytes, err := io.ReadAll(os.Stdin)
				if err != nil {
					return fmt.Errorf("failed to read from stdin: %w", err)
				}
				content = string(bytes)
			} else {
				content = strings.Join(args, " ")
			}

			if strings.TrimSpace(content) == "" {
				return fmt.Errorf("no content provided")
			}

			// Get project root
			root, err := getProjectRoot()
			if err != nil {
				return fmt.Errorf("failed to get project root: %w", err)
			}

			// Get sprint file path
			sprintFile := filepath.Join(root, "sprint.current.md")
			
			// Check if file exists
			if _, err := os.Stat(sprintFile); os.IsNotExist(err) {
				return fmt.Errorf("sprint.current.md not found. Use 'init' to create a new sprint")
			}

			// Format content with timestamp and separator
			formattedContent := "\n\n" + strings.Repeat("-", 80) + "\n\n" + formatSprintEntry(content)

			// Append to file
			f, err := os.OpenFile(sprintFile, os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				return fmt.Errorf("failed to open sprint file: %w", err)
			}
			defer f.Close()

			if _, err := f.WriteString(formattedContent); err != nil {
				return fmt.Errorf("failed to update sprint file: %w", err)
			}

			fmt.Printf("✨ Updated sprint.current.md in %s\n", root)
			return nil
		},
	}
}

// formatSprintEntry formats a sprint entry with timestamp
func formatSprintEntry(content string) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("## Sprint Update - %s\n\n%s", timestamp, content)
}

// getProjectRoot returns the root directory of the project
func getProjectRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to find project root: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}
