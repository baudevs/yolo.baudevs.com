package commands

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/baudevs/yolo.baudevs.com/internal/ai"
	"github.com/baudevs/yolo.baudevs.com/internal/license"
	"github.com/spf13/cobra"
)

// NewAskCommand initializes the ask command
func NewAskCommand() *cobra.Command {
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
	config, err := ai.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load AI config: %w", err)
	}

	// Initialize license manager
	licenseManager, err := license.NewManager(license.Config{
		StripeSecretKey:  os.Getenv("STRIPE_SECRET_KEY"),
		DefaultOpenAIKey: os.Getenv("OPENAI_API_KEY"),
	})
	if err != nil {
		return fmt.Errorf("failed to initialize license manager: %w", err)
	}

	// Create AI client
	aiClient, err := ai.NewClient(config, licenseManager)
	if err != nil {
		return fmt.Errorf("failed to initialize AI client: %w", err)
	}

	// Get response from AI
	response, err := aiClient.Ask(context.Background(), question)
	if err != nil {
		return fmt.Errorf("failed to get AI response: %w", err)
	}

	// Print the response
	fmt.Println(strings.TrimSpace(response))
	return nil
}
