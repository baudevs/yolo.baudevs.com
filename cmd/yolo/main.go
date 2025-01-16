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
		Short: "🚀 YOLO - Your friendly project companion!",
		Long: `🌟 Welcome to YOLO - Your Awesome Project Companion! 🌟

Think of YOLO as your friendly project assistant that helps you:
✨ Keep track of your amazing ideas (we call them "epics")
📝 Break down big tasks into manageable pieces
🤖 Work smoothly with AI to make coding fun
🎨 See your project come to life in beautiful 3D
📚 Never lose track of what you've done

Perfect for:
👩‍💼 Product Managers tracking features
👨‍💻 Developers (both new and experienced!)
🎓 Students learning to code with AI
👔 Executives wanting project insights
🎯 Anyone with a cool project idea

No technical background needed - we'll guide you every step of the way! 🌈`,
	}

	// Add commands with fun descriptions
	rootCmd.AddCommand(commands.InitCmd())
	rootCmd.AddCommand(commands.EpicCmd())
	rootCmd.AddCommand(commands.TaskCmd())
	rootCmd.AddCommand(commands.FeatureCmd())
	rootCmd.AddCommand(commands.StatusCmd())
	rootCmd.AddCommand(commands.PromptCmd())
	rootCmd.AddCommand(commands.GraphCmd())
	rootCmd.AddCommand(commands.CommitCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("🚨 Oops! Something went wrong:", err)
		fmt.Println("💡 Don't worry! Try running 'yolo help' for friendly guidance!")
		os.Exit(1)
	}
} 