package commands

import (
	"github.com/spf13/cobra"
	"fmt"
)

func EpicCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "epic [name]",
		Short: "✨ Create a big, exciting idea!",
		Long: `🌟 Let's capture your amazing project vision! 

An "epic" is like a big dream or goal for your project. Think of it as:
🎯 A major feature you want to build
🎨 A big problem you want to solve
📚 A collection of related smaller tasks

For example:
- "Create an awesome user dashboard"
- "Make our app super fast"
- "Add cool AI features"

Don't worry about the details yet - we'll break it down into smaller pieces later!

Examples:
  yolo epic "Amazing User Dashboard"
  yolo epic "Super Speed Boost"
  yolo epic "AI Magic Features"`,
		RunE: runEpic,
	}

	cmd.Flags().StringP("description", "d", "", "📝 Tell us more about your idea!")
	cmd.Flags().StringP("status", "s", "planning", "🎯 Where are you at? (planning, in-progress, done)")
	
	return cmd
}

func runEpic(cmd *cobra.Command, args []string) error {
	fmt.Println("🌟 Creating your epic adventure...")
	
	// ... rest of the implementation ...
	
	fmt.Println("\n✨ Epic idea captured successfully!")
	fmt.Println("\n💡 What's next?")
	fmt.Println("1. Add features with 'yolo feature'")
	fmt.Println("2. Break it down into tasks with 'yolo task'")
	fmt.Println("3. See it in 3D with 'yolo graph'")
	
	return nil
} 