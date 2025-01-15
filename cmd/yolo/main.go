package main

import (
	"fmt"
	"os"

	"github.com/baudevs/yolo-cli/internal/commands"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "yolo",
		Short: "YOLO - You Observe, Log, and Oversee methodology CLI tool",
		Long: `YOLO is a comprehensive project management methodology designed for 
working effectively with LLM developers while maintaining complete project history.`,
	}

	// Add commands
	rootCmd.AddCommand(commands.InitCmd())
	rootCmd.AddCommand(commands.EpicCmd())
	rootCmd.AddCommand(commands.TaskCmd())
	rootCmd.AddCommand(commands.FeatureCmd())
	rootCmd.AddCommand(commands.StatusCmd())
	rootCmd.AddCommand(commands.PromptCmd())
	rootCmd.AddCommand(commands.GraphCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
} 