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

func FeatureCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "feature [description]",
		Short: " Create a focused feature",
		Long: `Create a new feature in your YOLO project.

A feature is a significant piece of functionality that delivers value to users. YOLO will:
1. Use AI to analyze your feature description
2. Find or suggest creating a parent epic
3. Generate suggested tasks to implement the feature
4. Create and link everything automatically

Examples:
  yolo feature "Add user authentication system"
  yolo feature "Implement real-time notifications"
  yolo feature "Create dashboard analytics"`,
		Args: cobra.ExactArgs(1),
		RunE: runFeature,
	}

	cmd.Flags().StringP("epic", "e", "", "Explicitly link to an epic (e.g., E001)")
	cmd.Flags().StringP("status", "s", "planning", "Feature status (planning, in-progress, done)")

	return cmd
}

func runFeature(cmd *cobra.Command, args []string) error {
	description := args[0]
	fmt.Println(" Creating your feature with AI...")

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

	// Load all work items
	items, err := relManager.LoadWorkItems(relationships.Epic, relationships.Feature, relationships.Task)
	if err != nil {
		return fmt.Errorf("failed to load work items: %w", err)
	}

	// Find or create parent epic
	var parentEpic *relationships.WorkItem
	var createNewEpic bool
	epicFlag, _ := cmd.Flags().GetString("epic")

	if epicFlag != "" {
		// Use specified epic
		for _, item := range items {
			if item.Type == relationships.Epic && item.ID == epicFlag {
				parentEpic = &item
				break
			}
		}
		if parentEpic == nil {
			return fmt.Errorf("specified epic %s not found", epicFlag)
		}
	} else {
		// Let AI suggest parent epic
		fmt.Println(" Finding the best epic match...")
		parentEpic, createNewEpic, err = relManager.FindOrCreateParent(context.Background(), relationships.Feature, description, items)
		if err != nil {
			return fmt.Errorf("failed to find parent epic: %w", err)
		}
	}

	// If we need to create a new epic, do it now
	if createNewEpic {
		fmt.Println(" Creating new parent epic...")
		epicPrompt := fmt.Sprintf(`Create an epic description for a feature described as:
"%s"

The epic should:
1. Be broader in scope than the feature
2. Provide strategic context
3. Allow for related features

Respond with a concise but comprehensive epic description.`, description)

		epicContent, err := client.Ask(context.Background(), epicPrompt)
		if err != nil {
			return fmt.Errorf("failed to generate epic content: %w", err)
		}

		epicID := utils.GenerateID("E")
		epicTitle := fmt.Sprintf("Epic for %s", description)
		epicPath := filepath.Join("yolo", "epics", fmt.Sprintf("%s.md", epicID))

		epicFileContent := fmt.Sprintf(`# [%s] %s

## Status: planning
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
`, epicID, epicTitle,
			time.Now().Format("2006-01-02"),
			time.Now().Format("2006-01-02"),
			strings.TrimSpace(epicContent))

		if err := os.MkdirAll(filepath.Dir(epicPath), 0755); err != nil {
			return fmt.Errorf("failed to create epics directory: %w", err)
		}

		if err := os.WriteFile(epicPath, []byte(epicFileContent), 0644); err != nil {
			return fmt.Errorf("failed to write epic file: %w", err)
		}

		parentEpic = &relationships.WorkItem{
			Type:        relationships.Epic,
			ID:          epicID,
			Title:       epicTitle,
			Description: epicContent,
			Status:      "planning",
			Path:        epicPath,
			Content:     epicFileContent,
		}
		items = append(items, *parentEpic)
	}

	// Generate feature content
	featurePrompt := fmt.Sprintf(`Create a detailed feature description for:
"%s"

Consider this epic's context:
%s

The description should include:
1. Specific functionality to be implemented
2. User value and benefits
3. Technical considerations
4. Success criteria`, description, parentEpic.Description)

	fmt.Println(" Generating detailed feature description...")
	featureContent, err := client.Ask(context.Background(), featurePrompt)
	if err != nil {
		return fmt.Errorf("failed to generate feature content: %w", err)
	}

	// Create feature file
	status, _ := cmd.Flags().GetString("status")
	featureID := utils.GenerateID("F")
	featurePath := filepath.Join("yolo", "features", fmt.Sprintf("%s.md", featureID))

	featureFileContent := fmt.Sprintf(`# [%s] %s

## Status: %s
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
`, featureID, description, status,
		time.Now().Format("2006-01-02"),
		time.Now().Format("2006-01-02"),
		parentEpic.ID, parentEpic.Title,
		strings.TrimSpace(featureContent),
		parentEpic.ID, parentEpic.Title)

	if err := os.MkdirAll(filepath.Dir(featurePath), 0755); err != nil {
		return fmt.Errorf("failed to create features directory: %w", err)
	}

	if err := os.WriteFile(featurePath, []byte(featureFileContent), 0644); err != nil {
		return fmt.Errorf("failed to write feature file: %w", err)
	}

	currentFeature := relationships.WorkItem{
		Type:        relationships.Feature,
		ID:          featureID,
		Title:       description,
		Description: featureContent,
		Status:      status,
		Path:        featurePath,
		Content:     featureFileContent,
	}
	items = append(items, currentFeature)

	// Generate and create tasks
	fmt.Println(" Generating implementation tasks...")
	taskTitles, err := relManager.SuggestChildren(context.Background(), relationships.Feature, description)
	if err != nil {
		return fmt.Errorf("failed to generate tasks: %w", err)
	}

	var tasks []relationships.WorkItem
	for _, taskTitle := range taskTitles {
		taskPrompt := fmt.Sprintf(`Create a detailed task description for:
"%s"

This task is part of the feature:
"%s"

The description should be specific, actionable, and include clear success criteria.`, taskTitle, description)

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
			featureID, description,
			parentEpic.ID, parentEpic.Title,
			strings.TrimSpace(taskContent),
			featureID, description,
			parentEpic.ID, parentEpic.Title)

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
		tasks = append(tasks, task)
		items = append(items, task)

		fmt.Printf(" Created task: %s - %s\n", taskID, taskTitle)
	}

	// Update relationships in all files
	fmt.Println(" Updating relationships...")

	// Update epic relationships
	epicRelations := map[relationships.WorkItemType][]relationships.WorkItem{
		relationships.Feature: {currentFeature},
		relationships.Task:    tasks,
	}
	if err := relManager.UpdateRelationships(parentEpic.Path, epicRelations); err != nil {
		return fmt.Errorf("failed to update epic relationships: %w", err)
	}

	// Update feature relationships
	featureRelations := map[relationships.WorkItemType][]relationships.WorkItem{
		relationships.Epic: {*parentEpic},
		relationships.Task: tasks,
	}
	if err := relManager.UpdateRelationships(featurePath, featureRelations); err != nil {
		return fmt.Errorf("failed to update feature relationships: %w", err)
	}

	fmt.Printf("\n Feature %s created successfully!\n", featureID)
	fmt.Printf(" Linked to epic: [%s] %s\n", parentEpic.ID, parentEpic.Title)
	fmt.Printf(" Created %d implementation tasks\n", len(tasks))

	return nil
}
