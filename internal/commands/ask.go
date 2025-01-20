package commands

import (
	"fmt"
	"strings"

	"github.com/baudevs/yolo.baudevs.com/internal/ai"
	"github.com/baudevs/yolo.baudevs.com/internal/messages"
	"github.com/spf13/cobra"
)

// AskCmd initializes the ask command
func AskCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ask [question]",
		Short: "Ask an AI programming question",
		Long: `Ask any programming question and get a witty response with 3 simple steps to solve it.
Examples:
  yolo ask "how do I center a div?"
  yolo ask "what's the best way to handle errors in Go?"`,
		Args: cobra.ExactArgs(1),
		RunE: runAsk,
	}

	return cmd
}

func runAsk(cmd *cobra.Command, args []string) error {
	question := args[0]

	// Load AI config
	aiConfig, err := ai.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load AI config: %w", err)
	}

	// Create AI client
	aiClient := ai.NewClient(aiConfig)

	// Get current personality level
	personality := messages.GetPersonality()

	// Create prompt based on personality
	var prompt string
	switch personality {
	case messages.NerdyClean:
		prompt = fmt.Sprintf(`You are a helpful programming assistant. Please provide a 3-step solution to this question:
%s

Format your response as:
"Here's how to solve that, fellow developer!

1. [First step]
2. [Second step]
3. [Third step]"

Keep each step concise and practical.`, question)

	case messages.MildlyRude:
		prompt = fmt.Sprintf(`You are a slightly sassy programming assistant. Please provide a 3-step solution to this question:
%s

Format your response as:
"*rolls eyes* Fine, here's how you do it:

1. [First step]
2. [Second step]
3. [Third step]"

Keep each step concise and add a bit of attitude.`, question)

	case messages.UnhingedFunny:
		prompt = fmt.Sprintf(`You are an unhinged but hilarious programming assistant. Please provide a 3-step solution to this question:
%s

Format your response as:
"HOLY SH*T, LET'S F*CKING DO THIS! 🚀

1. [First step]
2. [Second step]
3. [Third step]"

Keep each step concise and add some wild energy!`, question)
	}

	// Get response from AI
	response, err := aiClient.GetCompletion(prompt)
	if err != nil {
		return fmt.Errorf("failed to get AI response: %w", err)
	}

	// Print the response
	fmt.Println(strings.TrimSpace(response))
	return nil
}
