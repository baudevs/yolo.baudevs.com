package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

// CommitAI handles AI-powered commit message generation
type CommitAI struct {
	client *openai.Client
}

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

// NewCommitAI creates a new CommitAI instance
func NewCommitAI() (*CommitAI, error) {
	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("no AI provider configured")
	}

	client := openai.NewClient(key)
	return &CommitAI{
		client: client,
	}, nil
}

// GenerateCommitMessage generates a commit message based on the changes
func (ai *CommitAI) GenerateCommitMessage(changes string) (string, error) {
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

Here are some examples of good commit messages:

Example 1 - New Feature:
{
  "type": "feat",
  "scope": "auth",
  "subject": "add OAuth2 authentication flow",
  "body": "Implement OAuth2 authentication with Google and GitHub providers:\n- Add OAuth routes and handlers\n- Create user session management\n- Add secure token storage\n- Implement callback handling",
  "breaking": false,
  "issue_refs": ["#123", "#124"],
  "co_authors": ["jane@example.com"]
}

Example 2 - Bug Fix:
{
  "type": "fix",
  "scope": "api",
  "subject": "handle null response in user profile endpoint",
  "body": "Add null checks to prevent crashes when user profile is incomplete",
  "breaking": false
}

Example 3 - Breaking Change:
{
  "type": "feat",
  "scope": "database",
  "subject": "migrate to PostgreSQL",
  "body": "Switch from SQLite to PostgreSQL for better scalability:\n- Update database schema\n- Migrate existing data\n- Update connection handling",
  "breaking": true,
  "issue_refs": ["#234"]
}

Example 4 - Documentation:
{
  "type": "docs",
  "scope": "readme",
  "subject": "update installation instructions",
  "body": "Add detailed steps for Windows and Linux installation"
}

Example 5 - Refactor:
{
  "type": "refactor",
  "scope": "core",
  "subject": "simplify error handling",
  "body": "Consolidate error handling into a central utility:\n- Create error types\n- Add context preservation\n- Improve error messages"
}

Now analyze these changes and generate a similar commit message:
` + changes

	resp, err := ai.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a helpful assistant that generates conventional commit messages in JSON format. Always ensure the output is valid JSON and follows the examples provided.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: openai.ChatCompletionResponseFormatTypeJSONObject,
			},
		},
	)
	if err != nil {
		return "", err
	}

	// Validate JSON structure
	var msg CommitMessage
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &msg); err != nil {
		return "", fmt.Errorf("invalid JSON response: %w", err)
	}

	return resp.Choices[0].Message.Content, nil
} 