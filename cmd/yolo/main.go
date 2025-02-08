package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/baudevs/yolo.baudevs.com/cmd"
	"github.com/baudevs/yolo.baudevs.com/internal/commands"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "yolo",
	Short: "🚀 YOLO - Your friendly project companion!",
	Long: `🌈 Welcome to YOLO - where project management meets fun! 

YOLO helps you:
✨ Create and organize projects
🤖 Get AI-powered assistance
📊 Track progress beautifully
🎮 See your work come alive

Perfect for:
👩‍💼 Product Managers
👨‍💻 Developers
🎨 Designers
📈 Team Leads
🎓 Students

No complicated stuff - just run a command and watch the magic happen!`,
}

func init() {
	// Load .env file from the project root
	if home, err := os.UserHomeDir(); err == nil {
		envPath := filepath.Join(home, ".yolo", ".env")
		_ = godotenv.Load(envPath)
	}

	// Also try loading from current directory
	_ = godotenv.Load()

	// Core commands
	rootCmd.AddCommand(commands.InitCmd())
	rootCmd.AddCommand(commands.ExplainCmd())
	rootCmd.AddCommand(cmd.NewDevCommand())
	rootCmd.AddCommand(commands.NewVersionCommand())

	// AI commands
	rootCmd.AddCommand(commands.NewAICommand())
	rootCmd.AddCommand(commands.NewCommitCommand())
	rootCmd.AddCommand(commands.NewAskCommand())

	// Prompt management
	rootCmd.AddCommand(commands.NewPromptCommand())
	rootCmd.AddCommand(commands.InitMessagePromptsCmd())

	// Project management
	rootCmd.AddCommand(commands.SprintCmd())
	rootCmd.AddCommand(commands.EpicCmd())
	rootCmd.AddCommand(commands.FeatureCmd())
	rootCmd.AddCommand(commands.TaskCmd())
	rootCmd.AddCommand(commands.GraphCmd) // Added Graph command

	// License management
	rootCmd.AddCommand(commands.NewLicenseCommand())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
