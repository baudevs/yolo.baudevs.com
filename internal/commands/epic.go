package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/baudevs/yolo.baudevs.com/internal/ai"
	"github.com/baudevs/yolo.baudevs.com/internal/config"
	"github.com/baudevs/yolo.baudevs.com/internal/license"
	"github.com/spf13/cobra"
)

func EpicCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "epic [description]",
		Short: "✨ Create a big, exciting idea!",
		Long: `🌟 Let's capture your amazing project vision! 

An "epic" is like a big dream or goal for your project. Think of it as:
🎯 A major feature you want to build
🎨 A big problem you want to solve
📚 A collection of related smaller tasks

Just provide a description of your epic, and YOLO will use AI to:
1. Generate detailed epic content
2. Break it down into 10 focused tasks
3. Create all necessary files and links

Examples:
  yolo epic "Create an AI-powered search feature"
  yolo epic "Implement real-time collaboration"
  yolo epic "Build a beautiful dashboard"`,
		Args: cobra.ExactArgs(1),
		RunE: runEpic,
	}

	cmd.Flags().StringP("status", "s", "planning", "🎯 Where are you at? (planning, in-progress, done)")
	
	return cmd
}

func runEpic(cmd *cobra.Command, args []string) error {
	description := args[0]
	fmt.Println("🌟 Creating your epic adventure with AI...")

	// Load config and create AI client
	clientConfig, err := config.LoadClientConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	licenseManager, err := license.NewManager()
	if err != nil {
		return fmt.Errorf("failed to create license manager: %w", err)
	}

	aiClient, err := ai.NewClient(clientConfig, licenseManager)
	if err != nil {
		return fmt.Errorf("failed to create AI client: %w", err)
	}

	// Generate epic content using AI
	epicPrompt := fmt.Sprintf(`Generate a detailed epic description and 10 tasks for the following project idea:
Description: %s

Respond in the following format:
[EPIC]
<Detailed epic description including goals, scope, and expected outcomes>

[TASKS]
1. <Task 1 title> | <Task 1 description>
2. <Task 2 title> | <Task 2 description>
...and so on for all 10 tasks

Each task should be specific, actionable, and contribute to the epic's completion.`, description)

	fmt.Println("🤖 Consulting AI for epic details and tasks...")
	epicContent, err := aiClient.Ask(context.Background(), epicPrompt)
	if err != nil {
		return fmt.Errorf("failed to generate epic content: %w", err)
	}

	// Parse AI response
	sections := strings.Split(epicContent, "[TASKS]")
	if len(sections) != 2 {
		return fmt.Errorf("invalid AI response format")
	}

	epicDescription := strings.TrimPrefix(sections[0], "[EPIC]\n")
	tasksContent := sections[1]

	// Create epic file
	status, _ := cmd.Flags().GetString("status")
	epicID := generateID("E")
	epicPath := filepath.Join("yolo", "epics", fmt.Sprintf("%s.md", epicID))

	epicFileContent := fmt.Sprintf(`# [%s] %s

## Status: %s
Created: %s
Last Updated: %s

## Description
%s

## Tasks
`, epicID, description, status, time.Now().Format("2006-01-02"), time.Now().Format("2006-01-02"), strings.TrimSpace(epicDescription))

	if err := os.MkdirAll(filepath.Dir(epicPath), 0755); err != nil {
		return fmt.Errorf("failed to create epics directory: %w", err)
	}

	if err := os.WriteFile(epicPath, []byte(epicFileContent), 0644); err != nil {
		return fmt.Errorf("failed to write epic file: %w", err)
	}

	// Create task files
	fmt.Println("\n📝 Creating tasks...")
	taskLines := strings.Split(tasksContent, "\n")
	for _, line := range taskLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Parse task number, title, and description
		parts := strings.SplitN(line, ".", 2)
		if len(parts) != 2 {
			continue
		}

		taskParts := strings.Split(parts[1], "|")
		if len(taskParts) != 2 {
			continue
		}

		taskTitle := strings.TrimSpace(taskParts[0])
		taskDescription := strings.TrimSpace(taskParts[1])

		// Create task file
		taskID := generateID("T")
		taskPath := filepath.Join("yolo", "tasks", fmt.Sprintf("%s.md", taskID))

		taskFileContent := fmt.Sprintf(`# [%s] %s

## Status: planning
Created: %s
Last Updated: %s
Epic: [%s] %s

## Description
%s
`, taskID, taskTitle, time.Now().Format("2006-01-02"), time.Now().Format("2006-01-02"), epicID, description, taskDescription)

		if err := os.MkdirAll(filepath.Dir(taskPath), 0755); err != nil {
			return fmt.Errorf("failed to create tasks directory: %w", err)
		}

		if err := os.WriteFile(taskPath, []byte(taskFileContent), 0644); err != nil {
			return fmt.Errorf("failed to write task file: %w", err)
		}

		fmt.Printf("✅ Created task: %s - %s\n", taskID, taskTitle)
	}

	fmt.Printf("\n✨ Epic %s created successfully!\n", epicID)
	fmt.Println("\n💡 What's next?")
	fmt.Println("1. Review the generated epic and tasks")
	fmt.Println("2. Start working on tasks with 'yolo task'")
	fmt.Println("3. See your progress in 3D with 'yolo graph'")

	return nil
}

func generateID(prefix string) string {
	// Get list of existing files
	files, err := filepath.Glob(filepath.Join("yolo", strings.ToLower(prefix)+"*", "*.md"))
	if err != nil {
		return fmt.Sprintf("%s001", prefix)
	}

	// Find highest number
	maxNum := 0
	for _, file := range files {
		base := filepath.Base(file)
		if len(base) < 4 {
			continue
		}
		numStr := strings.TrimPrefix(strings.TrimSuffix(base, ".md"), prefix)
		num := 0
		fmt.Sscanf(numStr, "%d", &num)
		if num > maxNum {
			maxNum = num
		}
	}

	return fmt.Sprintf("%s%03d", prefix, maxNum+1)
}