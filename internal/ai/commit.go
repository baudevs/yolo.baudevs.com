package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"github.com/baudevs/yolo.baudevs.com/internal/models"
	"github.com/sashabaranov/go-openai"
)

// CommitAI handles AI-powered commit message generation
type CommitAI struct {
	client *openai.Client
}

// NewCommitAI creates a new CommitAI instance
func NewCommitAI(apiKey string) (*CommitAI, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("OpenAI API key is required")
	}

	return &CommitAI{
		client: openai.NewClient(apiKey),
	}, nil
}

// Debug logs information to both console and file
func (ai *CommitAI) debug(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println("DEBUG:", msg)
	
	debugDir := "/tmp/yolo-debug"
	if err := os.MkdirAll(debugDir, 0755); err != nil {
		fmt.Printf("Warning: Could not create debug directory: %v\n", err)
		return
	}

	logFile := filepath.Join(debugDir, "yolo-commit.log")
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Warning: Could not open log file: %v\n", err)
		return
	}
	defer f.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(f, "[%s] %s\n", timestamp, msg)
}

// analyzeChunk sends a portion of changes to OpenAI and returns the analysis
func (ai *CommitAI) analyzeChunk(changes string, chunkNum, totalChunks int) (models.CommitMessage, error) {
	prompt := fmt.Sprintf(`Analyze the following Git changes (part %d of %d) and generate a conventional commit message in JSON format.
Note: Focus on understanding the changes in this chunk, a final summary will be generated later.

Follow these rules:
1. Use semantic commit types: feat, fix, docs, style, refactor, perf, test, build, ci, chore
2. Keep the subject line clear and concise
3. Use present tense ("add" not "added")
4. Focus on the main changes in this chunk
5. Include any breaking changes found

The response must be a valid JSON object with this structure:
{
  "type": "feat|fix|docs|style|refactor|perf|test|build|ci|chore",
  "scope": "optional area affected",
  "subject": "concise description",
  "body": "key changes in this chunk",
  "breaking": false,
  "issue_refs": ["optional array of issue references"],
  "co_authors": ["optional array of co-authors"]
}

Changes to analyze:
%s`, chunkNum, totalChunks, changes)

	ai.debug("Analyzing chunk %d of %d (length: %d bytes)", chunkNum, totalChunks, len(changes))

	resp, err := ai.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.2,
		},
	)

	if err != nil {
		return models.CommitMessage{}, fmt.Errorf("failed to analyze chunk %d: %w", chunkNum, err)
	}

	if len(resp.Choices) == 0 {
		return models.CommitMessage{}, fmt.Errorf("no analysis generated for chunk %d", chunkNum)
	}

	content := strings.TrimSpace(resp.Choices[0].Message.Content)
	var msg models.CommitMessage
	if err := json.Unmarshal([]byte(content), &msg); err != nil {
		return models.CommitMessage{}, fmt.Errorf("failed to parse chunk %d analysis: %w", chunkNum, err)
	}

	return msg, nil
}

// summarizeAnalyses combines multiple chunk analyses into a final commit message
func (ai *CommitAI) summarizeAnalyses(analyses []models.CommitMessage) (models.CommitMessage, error) {
	// Convert analyses to a JSON array for the AI to process
	analysesJSON, err := json.MarshalIndent(analyses, "", "  ")
	if err != nil {
		return models.CommitMessage{}, fmt.Errorf("failed to marshal analyses: %w", err)
	}

	prompt := fmt.Sprintf(`Analyze these commit message summaries and create a single, comprehensive commit message.
The summaries represent different parts of a large change set.

Previous analyses:
%s

Generate a final commit message that:
1. Uses the most appropriate commit type based on all changes
2. Creates a clear, concise subject line that captures the main change
3. Includes a body that summarizes the key changes
4. Preserves any breaking changes, issue references, or co-authors
5. Follows conventional commit format

Respond with a single JSON object using the same structure as the input.`, string(analysesJSON))

	resp, err := ai.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.2,
		},
	)

	if err != nil {
		return models.CommitMessage{}, fmt.Errorf("failed to generate final summary: %w", err)
	}

	if len(resp.Choices) == 0 {
		return models.CommitMessage{}, fmt.Errorf("no final summary generated")
	}

	content := strings.TrimSpace(resp.Choices[0].Message.Content)
	var finalMsg models.CommitMessage
	if err := json.Unmarshal([]byte(content), &finalMsg); err != nil {
		return models.CommitMessage{}, fmt.Errorf("failed to parse final summary: %w", err)
	}

	return finalMsg, nil
}

// GenerateCommitMessage generates a commit message based on the changes
func (ai *CommitAI) GenerateCommitMessage(changes string) (string, models.CommitMessage, error) {
	ai.debug("Starting commit message generation")
	ai.debug("Total changes length: %d bytes", len(changes))

	// Save full changes to file for reference
	debugDir := "/tmp/yolo-debug"
	timestamp := time.Now().Unix()
	changesFile := filepath.Join(debugDir, fmt.Sprintf("changes-%d.txt", timestamp))
	if err := os.WriteFile(changesFile, []byte(changes), 0644); err != nil {
		ai.debug("Warning: Could not write changes file: %v", err)
	}

	// Split changes into smaller chunks
	const maxChunkSize = 6000 // Conservative size to stay under token limit
	var chunks []string
	lines := strings.Split(changes, "\n")
	
	currentChunk := strings.Builder{}
	currentSize := 0
	
	for _, line := range lines {
		lineSize := len(line) + 1 // +1 for newline
		if currentSize+lineSize > maxChunkSize && currentSize > 0 {
			chunks = append(chunks, currentChunk.String())
			currentChunk.Reset()
			currentSize = 0
		}
		currentChunk.WriteString(line)
		currentChunk.WriteString("\n")
		currentSize += lineSize
	}
	
	if currentChunk.Len() > 0 {
		chunks = append(chunks, currentChunk.String())
	}

	ai.debug("Split changes into %d chunks", len(chunks))

	// Analyze each chunk
	var analyses []models.CommitMessage
	for i, chunk := range chunks {
		chunkNum := i + 1
		analysis, err := ai.analyzeChunk(chunk, chunkNum, len(chunks))
		if err != nil {
			return "", models.CommitMessage{}, fmt.Errorf("failed to analyze chunk %d: %w", chunkNum, err)
		}
		analyses = append(analyses, analysis)
		ai.debug("Successfully analyzed chunk %d", chunkNum)
	}

	// Generate final summary
	if len(analyses) == 1 {
		// If only one chunk, use its analysis directly
		formattedMsg := ai.FormatCommitMessage(analyses[0])
		ai.debug("Generated commit message from single chunk: %s", formattedMsg)
		return formattedMsg, analyses[0], nil
	}

	// Combine analyses into final message
	finalMsg, err := ai.summarizeAnalyses(analyses)
	if err != nil {
		return "", models.CommitMessage{}, fmt.Errorf("failed to generate final summary: %w", err)
	}

	formattedMsg := ai.FormatCommitMessage(finalMsg)
	ai.debug("Generated final commit message: %s", formattedMsg)
	return formattedMsg, finalMsg, nil
}

// FormatCommitMessage formats a CommitMessage into a conventional commit string
func (ai *CommitAI) FormatCommitMessage(msg models.CommitMessage) string {
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