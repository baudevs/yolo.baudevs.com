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

func TaskCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "task [description]",
		Short: " Create a focused task",
		Long: `Create a new task in your YOLO project.

A task is a small, well-defined unit of work. YOLO will:
1. Use AI to analyze your task description
2. Find the most relevant epic to link it to
3. Create the task with proper formatting and linking

Examples:
  yolo task "Add user authentication to the login page"
  yolo task "Implement dark mode in the UI"
  yolo task "Fix performance issues in the search feature"`,
		Args: cobra.ExactArgs(1),
		RunE: runTask,
	}

	cmd.Flags().StringP("epic", "e", "", "Explicitly link to an epic (e.g., E001)")
	cmd.Flags().StringP("status", "s", "planning", "Task status (planning, in-progress, done)")

	return cmd
}

func runTask(cmd *cobra.Command, args []string) error {
	description := args[0]
	fmt.Println(" Creating your task with AI...")

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
		parentEpic, createNewEpic, err = relManager.FindOrCreateParent(context.Background(), relationships.Task, description, items)
		if err != nil {
			return fmt.Errorf("failed to find parent epic: %w", err)
		}
	}

	// If we need to create a new epic, do it now
	if createNewEpic {
		fmt.Println(" Creating new parent epic...")
		epicPrompt := fmt.Sprintf(`Create an epic description for a task described as:
"%s"

The epic should:
1. Be broader in scope than the task
2. Provide implementation context
3. Allow for related tasks

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
- [ ] Epic implemented
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

	// Generate task content
	taskPrompt := fmt.Sprintf(`Create a detailed task description for:
"%s"

Consider this context:
%s

The description should be specific, actionable, and include clear success criteria.`, description, parentEpic.Description)

	fmt.Println(" Generating detailed task description...")
	taskContent, err := client.Ask(context.Background(), taskPrompt)
	if err != nil {
		return fmt.Errorf("failed to generate task content: %w", err)
	}

	// Create task file
	status, _ := cmd.Flags().GetString("status")
	taskID := utils.GenerateID("T")
	taskPath := filepath.Join("yolo", "tasks", fmt.Sprintf("%s.md", taskID))

	taskFileContent := fmt.Sprintf(`# [%s] %s

## Status: %s
Created: %s
Last Updated: %s
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
- Parent Epic: [%s] %s
<!-- YOLO-LINKS-END -->
`, taskID, description, status,
		time.Now().Format("2006-01-02"),
		time.Now().Format("2006-01-02"),
		parentEpic.ID, parentEpic.Title,
		strings.TrimSpace(taskContent),
		parentEpic.ID, parentEpic.Title)

	if err := os.MkdirAll(filepath.Dir(taskPath), 0755); err != nil {
		return fmt.Errorf("failed to create tasks directory: %w", err)
	}

	if err := os.WriteFile(taskPath, []byte(taskFileContent), 0644); err != nil {
		return fmt.Errorf("failed to write task file: %w", err)
	}

	currentTask := relationships.WorkItem{
		Type:        relationships.Task,
		ID:          taskID,
		Title:       description,
		Description: taskContent,
		Status:      status,
		Path:        taskPath,
		Content:     taskFileContent,
	}

	// Update relationships in all files
	fmt.Println(" Updating relationships...")

	// Update epic relationships
	epicRelations := map[relationships.WorkItemType][]relationships.WorkItem{
		relationships.Task: {currentTask},
	}
	if err := relManager.UpdateRelationships(parentEpic.Path, epicRelations); err != nil {
		return fmt.Errorf("failed to update epic relationships: %w", err)
	}

	fmt.Printf("\n Task %s created successfully!\n", taskID)
	fmt.Printf(" Linked to epic: [%s] %s\n", parentEpic.ID, parentEpic.Title)

	return nil
}
