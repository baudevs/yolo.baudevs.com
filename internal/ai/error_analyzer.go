package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type ErrorAnalysis struct {
	Problem     string   `json:"problem"`
	Explanation string   `json:"explanation"`
	Solutions   []string `json:"solutions"`
}

type ErrorAnalyzer struct {
	client *openai.Client
}

func NewErrorAnalyzer(apiKey string) *ErrorAnalyzer {
	return &ErrorAnalyzer{
		client: openai.NewClient(apiKey),
	}
}

func (ea *ErrorAnalyzer) AnalyzeError(err error, contextStr string) (*ErrorAnalysis, error) {
	prompt := fmt.Sprintf(`Analyze the following error in the context of a Git operation and provide a helpful explanation and solutions.
Context: %s
Error: %v

Respond with a JSON object containing:
{
  "problem": "brief description of the issue",
  "explanation": "user-friendly explanation of what went wrong",
  "solutions": ["array of step-by-step solutions"]
}`, contextStr, err)

	resp, err := ea.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to analyze error: %w", err)
	}

	var analysis ErrorAnalysis
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &analysis); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	return &analysis, nil
}

func (ea *ErrorAnalyzer) FormatAnalysis(analysis *ErrorAnalysis) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(" Problem: %s\n\n", analysis.Problem))
	sb.WriteString(fmt.Sprintf(" Explanation: %s\n\n", analysis.Explanation))
	sb.WriteString(" Solutions:\n")
	for i, solution := range analysis.Solutions {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, solution))
	}

	return sb.String()
}
