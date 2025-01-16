package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var CommitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Create a commit with AI-assisted message generation",
	Long: `Create a git commit with an AI-generated conventional commit message.
The command will analyze changes, generate a commit message using the configured AI provider,
and update project documentation automatically.`,
	RunE: runCommit,
}

func runCommit(cmd *cobra.Command, args []string) error {
	fmt.Println("🔍 Analyzing changes...")
	
	// Get git status
	changes, err := getGitChanges()
	if err != nil {
		return fmt.Errorf("failed to get git changes: %w", err)
	}

	if changes == "" {
		fmt.Println("✨ No changes to commit!")
		return nil
	}

	// Initialize AI provider
	ai, err := initAIProvider()
	if err != nil {
		return handleAIError(err)
	}

	// Generate commit message
	fmt.Println("🤖 Generating commit message...")
	message, err := ai.GenerateCommitMessage(changes)
	if err != nil {
		return handleCommitError(err)
	}

	// Stage changes
	fmt.Println("📦 Staging changes...")
	if err := stageChanges(); err != nil {
		return fmt.Errorf("failed to stage changes: %w", err)
	}

	// Create commit
	fmt.Println("💾 Creating commit...")
	if err := createCommit(message); err != nil {
		return fmt.Errorf("failed to create commit: %w", err)
	}

	// Update YOLO documentation
	fmt.Println("📝 Updating documentation...")
	if err := updateDocs(message); err != nil {
		return fmt.Errorf("failed to update documentation: %w", err)
	}

	// Stage documentation changes
	if err := stageChanges(); err != nil {
		return fmt.Errorf("failed to stage documentation: %w", err)
	}

	// Create documentation commit
	docMessage := fmt.Sprintf("docs: update YOLO documentation\n\n%s", message)
	if err := createCommit(docMessage); err != nil {
		return fmt.Errorf("failed to commit documentation: %w", err)
	}

	fmt.Println("✅ Commit created successfully!")
	return nil
}

func getGitChanges() (string, error) {
	cmd := exec.Command("git", "diff", "--staged")
	stagedOutput, err := cmd.Output()
	if err != nil {
		return "", err
	}

	cmd = exec.Command("git", "diff")
	unstagedOutput, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(stagedOutput) + string(unstagedOutput), nil
}

func stageChanges() error {
	cmd := exec.Command("git", "add", ".")
	return cmd.Run()
}

func createCommit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	return cmd.Run()
}

func updateDocs(message string) error {
	// Parse commit message for type, scope, and description
	parts := strings.SplitN(message, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid commit message format")
	}

	typeScope := strings.TrimSpace(parts[0])
	description := strings.TrimSpace(parts[1])

	// Extract type and scope
	typeParts := strings.Split(typeScope, "(")
	commitType := typeParts[0]
	scope := ""
	if len(typeParts) > 1 {
		scope = strings.TrimRight(typeParts[1], ")")
	}

	// Update HISTORY.yml
	entry := fmt.Sprintf(`
  - type: %s
    scope: %s
    subject: %q
    body: %q
`, commitType, scope, description, message)

	historyFile := "HISTORY.yml"
	history, err := os.OpenFile(historyFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer history.Close()

	if _, err := history.WriteString(entry); err != nil {
		return err
	}

	// Update CHANGELOG.md
	date := time.Now().Format("2006-01-02")
	changelogEntry := fmt.Sprintf("\n### %s\n- %s: %s", date, commitType, description)

	changelogFile := "CHANGELOG.md"
	changelog, err := os.OpenFile(changelogFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer changelog.Close()

	_, err = changelog.WriteString(changelogEntry)
	return err
}

func handleAIError(err error) error {
	fmt.Println("❌ AI provider error:")
	fmt.Println("You can configure an AI provider with the following steps:")
	fmt.Println("\n1. OpenAI API:")
	fmt.Println("   - Visit: https://platform.openai.com/api-keys")
	fmt.Println("   - Create a new API key")
	fmt.Println("   - Set environment variable: export OPENAI_API_KEY=your_key")
	
	fmt.Println("\n2. Anthropic Claude API:")
	fmt.Println("   - Visit: https://console.anthropic.com/")
	fmt.Println("   - Get an API key")
	fmt.Println("   - Set environment variable: export ANTHROPIC_API_KEY=your_key")
	
	fmt.Println("\n3. Mistral API:")
	fmt.Println("   - Visit: https://mistral.ai/")
	fmt.Println("   - Get an API key")
	fmt.Println("   - Set environment variable: export MISTRAL_API_KEY=your_key")
	
	fmt.Println("\nOr you can manually write your commit message:")
	fmt.Println("git commit -m \"type(scope): description\"")
	
	return fmt.Errorf("AI provider not configured: %w", err)
}

func handleCommitError(err error) error {
	fmt.Println("❌ Failed to generate commit message")
	fmt.Println("You can:")
	fmt.Println("1. Try again")
	fmt.Println("2. Write the message manually:")
	fmt.Println("   git commit -m \"type(scope): description\"")
	return err
} 