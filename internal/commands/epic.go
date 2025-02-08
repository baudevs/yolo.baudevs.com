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
	"github.com/baudevs/yolo.baudevs.com/internal/relationships"
	"github.com/baudevs/yolo.baudevs.com/internal/utils"
	"github.com/spf13/cobra"
)

func EpicCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "epic [description]",
		Short: "🌟 Create an epic",
		Long: `Create a new epic in your YOLO project.

An epic is a large body of work that can be broken down into features and tasks. YOLO will:
1. Use AI to analyze your epic description
2. Generate suggested features to implement the epic
3. Create implementation tasks for each feature
4. Link everything together automatically

Examples:
  yolo epic "Build user authentication system"
  yolo epic "Create analytics dashboard"
  yolo epic "Implement real-time collaboration"`,
		Args: cobra.ExactArgs(1),
		RunE: runEpic,
	}

	cmd.Flags().StringP("status", "s", "planning", "Epic status (planning, in-progress, done)")

	return cmd
}

func runEpic(cmd *cobra.Command, args []string) error {
	description := args[0]
	fmt.Println("🌟 Creating your epic with AI...")

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create license manager
	licenseManager, err := license.NewManager()
	if err != nil {
		return fmt.Errorf("failed to create license manager: %w", err)
	}

	// Create AI client
	client, err := ai.NewClient(cfg, licenseManager)
	if err != nil {
		return fmt.Errorf("failed to create AI client: %w", err)
	}

	// Create relationship manager
	relManager := relationships.NewManager(client)

	// Generate epic content with AI
	epicPrompt := fmt.Sprintf(`Create a comprehensive epic description for:
"%s"

The description should include:
1. Strategic goals and objectives
2. Business value and impact
3. High-level technical considerations
4. Success criteria and metrics

Make it detailed but concise.`, description)

	fmt.Println("🤖 Generating epic description...")
	epicContent, err := client.Ask(context.Background(), epicPrompt)
	if err != nil {
		return fmt.Errorf("failed to generate epic content: %w", err)
	}

	// Create epic file
	status, _ := cmd.Flags().GetString("status")
	epicID := utils.GenerateID("E")
	epicPath := filepath.Join("yolo", "epics", fmt.Sprintf("%s.md", epicID))

	epicFileContent := fmt.Sprintf(`# [%s] %s

## Status: %s
Created: %s
Last Updated: %s

## Description
%s

## Success Criteria
- [ ] Features implemented
- [ ] Tests added
- [ ] Documentation updated
- [ ] Code reviewed

## Relationships
<!-- YOLO-LINKS-START -->
<!-- YOLO-LINKS-END -->
`, epicID, description, status,
		time.Now().Format("2006-01-02"),
		time.Now().Format("2006-01-02"),
		strings.TrimSpace(epicContent))

	if err := os.MkdirAll(filepath.Dir(epicPath), 0755); err != nil {
		return fmt.Errorf("failed to create epics directory: %w", err)
	}

	if err := os.WriteFile(epicPath, []byte(epicFileContent), 0644); err != nil {
		return fmt.Errorf("failed to write epic file: %w", err)
	}

	currentEpic := relationships.WorkItem{
		Type:        relationships.Epic,
		ID:          epicID,
		Title:       description,
		Description: epicContent,
		Status:      status,
		Path:        epicPath,
		Content:     epicFileContent,
	}

	// Generate and create features
	fmt.Println("🤖 Generating implementation features...")
	featureTitles, err := relManager.SuggestChildren(context.Background(), relationships.Epic, description)
	if err != nil {
		return fmt.Errorf("failed to generate features: %w", err)
	}

	var features []relationships.WorkItem
	var allTasks []relationships.WorkItem

	for _, featureTitle := range featureTitles {
		featurePrompt := fmt.Sprintf(`Create a detailed feature description for:
"%s"

This feature is part of the epic:
"%s"

The description should include:
1. Specific functionality to be implemented
2. User value and benefits
3. Technical considerations
4. Success criteria`, featureTitle, description)

		featureContent, err := client.Ask(context.Background(), featurePrompt)
		if err != nil {
			return fmt.Errorf("failed to generate feature content: %w", err)
		}

		featureID := utils.GenerateID("F")
		featurePath := filepath.Join("yolo", "features", fmt.Sprintf("%s.md", featureID))

		featureFileContent := fmt.Sprintf(`# [%s] %s

## Status: planning
Created: %s
Last Updated: %s
Epic: [%s] %s

## Description
%s

## Success Criteria
- [ ] Feature implemented
- [ ] Tests added
- [ ] Documentation updated
- [ ] Code reviewed

## Relationships
<!-- YOLO-LINKS-START -->
- Parent Epic: [%s] %s
<!-- YOLO-LINKS-END -->
`, featureID, featureTitle,
			time.Now().Format("2006-01-02"),
			time.Now().Format("2006-01-02"),
			epicID, description,
			strings.TrimSpace(featureContent),
			epicID, description)

		if err := os.MkdirAll(filepath.Dir(featurePath), 0755); err != nil {
			return fmt.Errorf("failed to create features directory: %w", err)
		}

		if err := os.WriteFile(featurePath, []byte(featureFileContent), 0644); err != nil {
			return fmt.Errorf("failed to write feature file: %w", err)
		}

		feature := relationships.WorkItem{
			Type:        relationships.Feature,
			ID:          featureID,
			Title:       featureTitle,
			Description: featureContent,
			Status:      "planning",
			Path:        featurePath,
			Content:     featureFileContent,
		}
		features = append(features, feature)

		fmt.Printf("✅ Created feature: %s - %s\n", featureID, featureTitle)

		// Generate tasks for this feature
		fmt.Printf("🤖 Generating tasks for feature: %s\n", featureTitle)
		taskTitles, err := relManager.SuggestChildren(context.Background(), relationships.Feature, featureTitle)
		if err != nil {
			return fmt.Errorf("failed to generate tasks: %w", err)
		}

		for _, taskTitle := range taskTitles {
			taskPrompt := fmt.Sprintf(`Create a detailed task description for:
"%s"

This task is part of the feature:
"%s"

The description should be specific, actionable, and include clear success criteria.`, taskTitle, featureTitle)

			taskContent, err := client.Ask(context.Background(), taskPrompt)
			if err != nil {
				return fmt.Errorf("failed to generate task content: %w", err)
			}

			taskID := utils.GenerateID("T")
			taskPath := filepath.Join("yolo", "tasks", fmt.Sprintf("%s.md", taskID))

			taskFileContent := fmt.Sprintf(`# [%s] %s

## Status: planning
Created: %s
Last Updated: %s
Feature: [%s] %s
Epic: [%s] %s

## Description
%s

## Success Criteria
- [ ] Task implemented
- [ ] Code reviewed
- [ ] Tests added
- [ ] Documentation updated

## Relationships
<!-- YOLO-LINKS-START -->
- Parent Feature: [%s] %s
- Parent Epic: [%s] %s
<!-- YOLO-LINKS-END -->
`, taskID, taskTitle,
				time.Now().Format("2006-01-02"),
				time.Now().Format("2006-01-02"),
				featureID, featureTitle,
				epicID, description,
				strings.TrimSpace(taskContent),
				featureID, featureTitle,
				epicID, description)

			if err := os.MkdirAll(filepath.Dir(taskPath), 0755); err != nil {
				return fmt.Errorf("failed to create tasks directory: %w", err)
			}

			if err := os.WriteFile(taskPath, []byte(taskFileContent), 0644); err != nil {
				return fmt.Errorf("failed to write task file: %w", err)
			}

			task := relationships.WorkItem{
				Type:        relationships.Task,
				ID:          taskID,
				Title:       taskTitle,
				Description: taskContent,
				Status:      "planning",
				Path:        taskPath,
				Content:     taskFileContent,
			}
			allTasks = append(allTasks, task)

			fmt.Printf("✅ Created task: %s - %s\n", taskID, taskTitle)
		}

		// Update feature relationships
		featureRelations := map[relationships.WorkItemType][]relationships.WorkItem{
			relationships.Epic: {currentEpic},
			relationships.Task: allTasks,
		}
		if err := relManager.UpdateRelationships(featurePath, featureRelations); err != nil {
			return fmt.Errorf("failed to update feature relationships: %w", err)
		}
	}

	// Update epic relationships
	epicRelations := map[relationships.WorkItemType][]relationships.WorkItem{
		relationships.Feature: features,
		relationships.Task:    allTasks,
	}
	if err := relManager.UpdateRelationships(epicPath, epicRelations); err != nil {
		return fmt.Errorf("failed to update epic relationships: %w", err)
	}

	fmt.Printf("\n✨ Epic %s created successfully!\n", epicID)
	fmt.Printf("📋 Created %d features and %d tasks\n", len(features), len(allTasks))
	fmt.Println("\nNext steps:")
	fmt.Println("1. Review the generated content")
	fmt.Println("2. Assign features and tasks to team members")
	fmt.Println("3. See your progress in 3D with 'yolo graph'")

	return nil
}
