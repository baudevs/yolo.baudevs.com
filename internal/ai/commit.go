package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"github.com/baudevs/yolo.baudevs.com/internal/models"
	"github.com/sashabaranov/go-openai"
)


// CommitAI handles AI-powered commit message generation
type CommitAI struct {
	client *Client
}


// NewCommitAI creates a new CommitAI instance
func NewCommitAI(apiKey string) (*CommitAI, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("OpenAI API key is required")
	}

	config := &Config{
		DefaultProvider: "openai",
		Providers: map[string]Provider{
			"openai": {
				APIKey: apiKey,
			},
		},
	}

	return &CommitAI{
		client: NewClient(config),
	}, nil
}

// GenerateCommitMessage generates a commit message based on the changes
func (ai *CommitAI) GenerateCommitMessage(changes string) (string, models.CommitMessage, error) {
    prompt := `Analyze the following Git changes and generate a conventional commit message in JSON format.

Follow these rules:
1. Use semantic commit types: feat, fix, docs, style, refactor, perf, test, build, ci, chore
2. Keep the subject line clear and concise (max 72 chars)
3. Use present tense ("add" not "added")
4. Don't end the subject line with a period
5. Add relevant details in the body
6. Mark breaking changes with breaking=true
7. Include issue references if found in the changes
8. Add co-authors if multiple contributors are detected

The response must be a valid JSON object with this structure:
{
  "type": "feat|fix|docs|style|refactor|perf|test|build|ci|chore",
  "scope": "optional area affected",
  "subject": "concise description",
  "body": "optional detailed explanation",
  "breaking": false,
  "issue_refs": ["optional array of issue references"],
  "co_authors": ["optional array of co-authors"]
}

Changes to analyze:
` + changes

    content, err := ai.client.GetCompletion(prompt)
    if err != nil {
        return "", models.CommitMessage{}, fmt.Errorf("failed to generate commit message: %w", err)
    }

    content = strings.TrimSpace(content)

    fmt.Printf("📝 Raw AI response:\n%s\n", content)

    // If the response is not wrapped in curly braces, wrap it
    if !strings.HasPrefix(content, "{") {
        content = "{" + content
        fmt.Println("⚠️ Added opening brace")
    }
    if !strings.HasSuffix(content, "}") {
        content = content + "}"
        fmt.Println("⚠️ Added closing brace")
    }

    var msg models.CommitMessage
    if err := json.Unmarshal([]byte(content), &msg); err != nil {
        return "", models.CommitMessage{}, fmt.Errorf("invalid AI response format: %w\nResponse: %s", err, content)
    }

    formattedMessage := FormatCommitMessage(msg)
    return formattedMessage, msg, nil
}

// FormatCommitMessage formats a CommitMessage into a conventional commit string
func FormatCommitMessage(msg models.CommitMessage) string {
	var parts []string

	// Add type
	parts = append(parts, msg.Type)

	// Add scope if present
	if msg.Scope != "" {
		parts = append(parts, fmt.Sprintf("(%s)", msg.Scope))
	}

	// Add breaking change marker
	if msg.Breaking {
		parts = append(parts, "!")
	}

	// Add subject
	commitMsg := fmt.Sprintf("%s: %s", strings.Join(parts, ""), msg.Subject)

	// Add body if present
	if msg.Body != "" {
		commitMsg += "\n\n" + msg.Body
	}

	// Add issue references
	if len(msg.IssueRefs) > 0 {
		commitMsg += "\n\nRefs: " + strings.Join(msg.IssueRefs, ", ")
	}

	// Add co-authors
	if len(msg.CoAuthors) > 0 {
		for _, author := range msg.CoAuthors {
			commitMsg += fmt.Sprintf("\n\nCo-authored-by: %s", author)
		}
	}

	return commitMsg
}