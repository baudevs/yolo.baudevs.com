package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/baudevs/yolo.baudevs.com/internal/ai"
	"github.com/baudevs/yolo.baudevs.com/internal/git"
	"github.com/spf13/cobra"
)

type CommitMessage struct {
	Type        string   `json:"type"`
	Scope       string   `json:"scope,omitempty"`
	Subject     string   `json:"subject"`
	Body        string   `json:"body,omitempty"`
	Breaking    bool     `json:"breaking,omitempty"`
	IssueRefs   []string `json:"issue_refs,omitempty"`
	CoAuthors   []string `json:"co_authors,omitempty"`
}

type CommitOptions struct {
	NoSync bool
	Force  bool
}

func CommitCmd() *cobra.Command {
	opts := &CommitOptions{}

	cmd := &cobra.Command{
		Use:   "commit",
		Short: "Create an AI-powered commit with your changes",
		Long: `Create a commit with AI-generated message, automatically stage changes,
and sync with remote repository. The AI will analyze your changes and:
1. Generate a conventional commit message
2. Stage and commit your changes
3. Sync with remote (if available)
4. Help fix any errors that occur`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCommit(opts)
		},
	}

	cmd.Flags().BoolVar(&opts.NoSync, "no-sync", false, "Don't sync with remote")
	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Force commit even with warnings")

	return cmd
}

func runCommit(opts *CommitOptions) error {
	// Get working directory
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Initialize Git operations
	gitOps := git.NewGitOps(wd)

	// Check for changes
	if !gitOps.HasChanges() {
		fmt.Println("✨ Nothing to commit - working tree clean")
		return nil
	}

	// Get changes for AI analysis
	fmt.Println("🔍 Analyzing your changes...")
	changes, err := gitOps.GetChanges()
	if err != nil {
		return handleError("Failed to get changes", err, "git_changes")
	}

	// Initialize AI providers
	commitAI, err := ai.NewCommitAI(os.Getenv("OPENAI_API_KEY"))
	if err != nil {
		return handleError("Failed to initialize AI", err, "ai_init")
	}

	errorAnalyzer := ai.NewErrorAnalyzer(os.Getenv("OPENAI_API_KEY"))

	// Generate commit message
	fmt.Println("🤖 Generating commit message...")
	message, err := commitAI.GenerateCommitMessage(changes)
	if err != nil {
		return handleError("Failed to generate commit message", err, "ai_message")
	}

	// Parse the JSON response
	var commitMsg CommitMessage
	if err := json.Unmarshal([]byte(message), &commitMsg); err != nil {
		return handleError("Failed to parse commit message", err, "ai_parse")
	}

	// Format the conventional commit message
	formattedMessage := formatCommitMessage(commitMsg)

	// Stage changes
	fmt.Println("📦 Staging changes...")
	if err := gitOps.StageAll(); err != nil {
		return handleError("Failed to stage changes", err, "git_stage")
	}

	// Create commit
	fmt.Println("💾 Creating commit...")
	if err := gitOps.Commit(formattedMessage); err != nil {
		return handleError("Failed to create commit", err, "git_commit")
	}

	// Sync with remote if needed
	if !opts.NoSync && gitOps.HasRemote() {
		fmt.Println("🔄 Syncing with remote...")
		
		// Pull first
		if err := gitOps.Pull(); err != nil {
			analysis, _ := errorAnalyzer.AnalyzeError(err, "pulling from remote")
			if analysis != nil {
				fmt.Println("\n❌ Pull failed!")
				fmt.Println(errorAnalyzer.FormatAnalysis(analysis))
				
				if !opts.Force {
					return fmt.Errorf("sync failed")
				}
			}
		}

		// Then push
		if err := gitOps.Push(); err != nil {
			analysis, _ := errorAnalyzer.AnalyzeError(err, "pushing to remote")
			if analysis != nil {
				fmt.Println("\n❌ Push failed!")
				fmt.Println(errorAnalyzer.FormatAnalysis(analysis))
				return fmt.Errorf("sync failed")
			}
		}
	}

	// Update YOLO documentation
	fmt.Println("📝 Updating the project story...")
	if err := updateDocs(formattedMessage); err != nil {
		return handleError("Failed to update docs", err, "docs_update")
	}

	// Stage documentation changes
	if err := gitOps.StageAll(); err != nil {
		return handleError("Failed to stage doc updates", err, "git_stage_docs")
	}

	// Create documentation commit
	docMessage := fmt.Sprintf("docs: update YOLO documentation\n\n%s", formattedMessage)
	if err := gitOps.Commit(docMessage); err != nil {
		return handleError("Failed to create doc commit", err, "git_commit_docs")
	}

	fmt.Println("\n✅ All done! Here's what happened:")
	fmt.Printf("1. Created commit: %s\n", formattedMessage)
	if !opts.NoSync && gitOps.HasRemote() {
		fmt.Println("2. Synced with remote repository")
	}

	return nil
}

func formatCommitMessage(msg CommitMessage) string {
	// Start with type and scope
	result := msg.Type
	if msg.Scope != "" {
		result += fmt.Sprintf("(%s)", msg.Scope)
	}
	if msg.Breaking {
		result += "!"
	}
	result += fmt.Sprintf(": %s", msg.Subject)

	// Add body if present
	if msg.Body != "" {
		result += fmt.Sprintf("\n\n%s", msg.Body)
	}

	// Add issue references if present
	if len(msg.IssueRefs) > 0 {
		result += fmt.Sprintf("\n\nRefs: %s", strings.Join(msg.IssueRefs, ", "))
	}

	// Add co-authors if present
	for _, author := range msg.CoAuthors {
		result += fmt.Sprintf("\n\nCo-authored-by: %s", author)
	}

	return result
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

func handleError(context string, err error, errorType string) error {
	// Initialize error analyzer
	errorAnalyzer := ai.NewErrorAnalyzer(os.Getenv("OPENAI_API_KEY"))
	
	// Get AI analysis of the error
	analysis, analyzeErr := errorAnalyzer.AnalyzeError(err, context)
	if analyzeErr != nil {
		return fmt.Errorf("%s: %w", context, err)
	}

	// Print the analysis
	fmt.Printf("\n❌ %s\n", context)
	fmt.Println(errorAnalyzer.FormatAnalysis(analysis))
	
	return fmt.Errorf("%s: %w", context, err)
}