package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/baudevs/yolo-cli/internal/ai"
	"github.com/spf13/cobra"
)

// CommitMessage represents the structured output from AI
type CommitMessage struct {
	Type        string   `json:"type"`
	Scope       string   `json:"scope,omitempty"`
	Subject     string   `json:"subject"`
	Body        string   `json:"body,omitempty"`
	Breaking    bool     `json:"breaking,omitempty"`
	IssueRefs   []string `json:"issue_refs,omitempty"`
	CoAuthors   []string `json:"co_authors,omitempty"`
}

var CommitCmd = &cobra.Command{
	Use:   "commit",
	Short: "🤖 Save your changes with AI help!",
	Long: `✨ Let's save your amazing work with style! 

This magical command helps you:
1. 🔍 Look at what you've changed
2. 🤖 Ask AI to write a perfect description
3. 📦 Package everything up nicely
4. 📝 Update all the important docs
5. ✨ Keep your project history beautiful

No more worrying about:
❌ "What should I write in the commit message?"
❌ "Am I following the right format?"
❌ "Did I forget to update something?"

Our AI friend will:
🎯 Look at your changes
🎨 Write a clear description
📚 Follow best practices
🔄 Keep everything in sync

Perfect for:
👩‍💼 Product updates
👨‍💻 Code changes
📝 Documentation
🎨 Design tweaks
🐛 Bug fixes

Just run 'yolo commit' and we'll handle the rest!`,
	RunE: runCommit,
}

func runCommit(cmd *cobra.Command, args []string) error {
	fmt.Println("🔍 Looking at your amazing changes...")
	
	// Get git status
	changes, err := getGitChanges()
	if err != nil {
		return fmt.Errorf("❌ Oops! Couldn't see your changes: %w", err)
	}

	if changes == "" {
		fmt.Println("✨ Nothing to save yet - make some changes first!")
		return nil
	}

	// Initialize AI provider
	fmt.Println("🤖 Waking up our AI friend...")
	ai, err := initAIProvider()
	if err != nil {
		return handleAIError(err)
	}

	// Generate commit message
	fmt.Println("🎨 Creating the perfect description...")
	message, err := ai.GenerateCommitMessage(changes)
	if err != nil {
		return handleCommitError(err)
	}

	// Parse the JSON response
	var commitMsg CommitMessage
	if err := json.Unmarshal([]byte(message), &commitMsg); err != nil {
		return fmt.Errorf("❌ Oops! Couldn't understand AI's response: %w", err)
	}

	// Format the conventional commit message
	formattedMessage := formatCommitMessage(commitMsg)

	// Stage changes
	fmt.Println("📦 Packaging up your changes...")
	if err := stageChanges(); err != nil {
		return fmt.Errorf("❌ Oops! Couldn't package your changes: %w", err)
	}

	// Create commit
	fmt.Println("💾 Saving your work...")
	if err := createCommit(formattedMessage); err != nil {
		return fmt.Errorf("❌ Oops! Couldn't save your changes: %w", err)
	}

	// Update YOLO documentation
	fmt.Println("📝 Updating the project story...")
	if err := updateDocs(formattedMessage); err != nil {
		return fmt.Errorf("❌ Oops! Couldn't update the docs: %w", err)
	}

	// Stage documentation changes
	if err := stageChanges(); err != nil {
		return fmt.Errorf("❌ Oops! Couldn't package the doc updates: %w", err)
	}

	// Create documentation commit
	docMessage := fmt.Sprintf("docs: update YOLO documentation\n\n%s", formattedMessage)
	if err := createCommit(docMessage); err != nil {
		return fmt.Errorf("❌ Oops! Couldn't save the doc updates: %w", err)
	}

	fmt.Println("\n🎉 All done! Your changes are safely saved!")
	fmt.Println("\n💡 What's next?")
	fmt.Println("1. Make more amazing changes")
	fmt.Println("2. Run 'yolo status' to see how things are going")
	fmt.Println("3. Check 'yolo graph' to see your progress!")
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
	fmt.Println("🤖 Our AI friend needs a little help!")
	fmt.Println("\nLet's get you set up with an AI assistant. You can choose:")
	
	fmt.Println("\n1. ✨ OpenAI (ChatGPT)")
	fmt.Println("   🌐 Visit: https://platform.openai.com/api-keys")
	fmt.Println("   🎯 Create a new key")
	fmt.Println("   💻 Run: export OPENAI_API_KEY=your_key")
	
	fmt.Println("\n2. 🔮 Anthropic Claude")
	fmt.Println("   🌐 Visit: https://console.anthropic.com/")
	fmt.Println("   🎯 Get your key")
	fmt.Println("   💻 Run: export ANTHROPIC_API_KEY=your_key")
	
	fmt.Println("\n3. 🌟 Mistral AI")
	fmt.Println("   🌐 Visit: https://mistral.ai/")
	fmt.Println("   🎯 Get your key")
	fmt.Println("   💻 Run: export MISTRAL_API_KEY=your_key")
	
	fmt.Println("\n💡 Or, you can write your message manually:")
	fmt.Println("git commit -m \"type(area): what you did\"")
	
	return fmt.Errorf("🔑 No AI helper configured: %w", err)
}

func handleCommitError(err error) error {
	fmt.Println("❌ Oops! Something went wrong with the message")
	fmt.Println("\n💡 You can:")
	fmt.Println("1. 🔄 Try again")
	fmt.Println("2. ✍️  Write it yourself:")
	fmt.Println("   git commit -m \"type(area): what you did\"")
	return err
}

func initAIProvider() (*ai.CommitAI, error) {
	ai, err := ai.NewCommitAI()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize AI provider: %w", err)
	}
	return ai, nil
} 