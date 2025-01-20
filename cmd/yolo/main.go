package main

import (
	"os"

	"github.com/baudevs/yolo-cli/internal/commands"
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
	rootCmd.AddCommand(commands.InitCmd())
	rootCmd.AddCommand(commands.PromptCmd())
	rootCmd.AddCommand(commands.GraphCmd)
	rootCmd.AddCommand(commands.CommitCmd())
	rootCmd.AddCommand(commands.ShortcutsCmd)
	rootCmd.AddCommand(commands.ExplainCmd())
	rootCmd.AddCommand(commands.AICmd())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}