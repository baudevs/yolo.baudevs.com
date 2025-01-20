package commands

import (
	//"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
	"github.com/baudevs/yolo.baudevs.com/internal/ai"
	"github.com/baudevs/yolo.baudevs.com/internal/git"
	"github.com/baudevs/yolo.baudevs.com/internal/models"
	"github.com/spf13/cobra"
)



func CommitCmd() *cobra.Command {
	opts := &models.CommitOptions{}

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

func runCommit(opts *models.CommitOptions) error {
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

	// Load AI configuration
	config, err := ai.LoadConfig()
	if err != nil {
		return handleError("Failed to load AI configuration", err, "ai_config")
	}

	// Get API key from configuration
	apiKey := config.GetAPIKey(config.DefaultProvider)
	if apiKey == "" {
		return fmt.Errorf("OpenAI API key not found. Please run 'yolo ai configure' first")
	}

	// Initialize AI providers
	commitAI, err := ai.NewCommitAI(apiKey)
	if err != nil {
		return handleError("Failed to initialize AI", err, "ai_init")
	}

	errorAnalyzer := ai.NewErrorAnalyzer(apiKey)

	// Generate commit message
	fmt.Println("🤖 Generating commit message...")
	message,structuredMsg, err := commitAI.GenerateCommitMessage(changes)
	if err != nil {
		return handleError("Failed to generate commit message", err, "ai_message")
	}
	// use the structured message response
	commitMsg := structuredMsg
	
	// Parse the JSON response
	/* if err := json.Unmarshal([]byte(message), &commitMsg); err != nil {
		fmt.Printf("Failed to parse JSON: %v\n", err)
		return handleError("Failed to parse commit message", err, "ai_parse")
	}  */

	// Stage changes
	fmt.Println("📦 Staging changes...")
	if err := gitOps.StageAll(); err != nil {
		return handleError("Failed to stage changes", err, "git_stage")
	}

	// Create commit
	fmt.Println("💾 Creating commit...")
	if err := gitOps.Commit(message); err != nil {
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
	if err := updateDocs(commitMsg); err != nil {
		return handleError("Failed to update docs", err, "docs_update")
	}

	// Stage documentation changes
	if err := gitOps.StageAll(); err != nil {
		return handleError("Failed to stage doc updates", err, "git_stage_docs")
	}

	// Create documentation commit
	docMessage := fmt.Sprintf("docs: update YOLO documentation\n\n%s", formatCommitMessage(commitMsg))
	if err := gitOps.Commit(docMessage); err != nil {
		return handleError("Failed to create doc commit", err, "git_commit_docs")
	}

	fmt.Println("\n✅ All done! Here's what happened:")
	fmt.Printf("1. Created commit: %s\n", formatCommitMessage(commitMsg))
	if !opts.NoSync && gitOps.HasRemote() {
		fmt.Println("2. Synced with remote repository")
	}

	return nil
}

func formatCommitMessage(msg models.CommitMessage) string {
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

func updateDocs(msg models.CommitMessage) error {
	// Format the conventional commit message
	formattedMessage := formatCommitMessage(msg)

	// Update HISTORY.yml
	entry := fmt.Sprintf(`
  - type: %s
    scope: %s
    subject: %q
    body: %q
`, msg.Type, msg.Scope, msg.Subject, formattedMessage)

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
	changelogEntry := fmt.Sprintf("\n### %s\n- %s: %s", date, msg.Type, msg.Subject)

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
	// Load AI configuration
	config, configErr := ai.LoadConfig()
	if configErr != nil {
		return fmt.Errorf("%s: %w (failed to load AI config: %v)", context, err, configErr)
	}

	// Get API key from configuration
	apiKey := config.GetAPIKey(config.DefaultProvider)
	if apiKey == "" {
		return fmt.Errorf("%s: %w (OpenAI API key not found, run 'yolo ai configure' first)", context, err)
	}

	// Initialize error analyzer
	errorAnalyzer := ai.NewErrorAnalyzer(apiKey)
	
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