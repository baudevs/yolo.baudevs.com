package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/baudevs/yolo.baudevs.com/internal/ai"
	"github.com/baudevs/yolo.baudevs.com/internal/config"
	"github.com/baudevs/yolo.baudevs.com/internal/license"
	"github.com/baudevs/yolo.baudevs.com/internal/utils"
	"github.com/spf13/cobra"
)

type Epic struct {
	ID          string
	Title       string
	Description string
	Status      string
	Path        string
	Content     string
}

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

	// Load existing epics
	epics, err := loadEpics()
	if err != nil {
		return fmt.Errorf("failed to load epics: %w", err)
	}

	var targetEpic *Epic
	epicFlag, _ := cmd.Flags().GetString("epic")
	
	if epicFlag != "" {
		// Use specified epic
		for _, epic := range epics {
			if epic.ID == epicFlag {
				targetEpic = &epic
				break
			}
		}
		if targetEpic == nil {
			return fmt.Errorf("specified epic %s not found", epicFlag)
		}
	} else {
		// Ask AI to find the best matching epic
		epicPrompt := fmt.Sprintf(`Given this task description: "%s"

Analyze these existing epics and determine which one is the most relevant to link the task to. Consider the epic's goals, scope, and existing features. If none are relevant, respond with "NONE".

Epics:
%s

Respond with just the epic ID (e.g., "E001") or "NONE" if no epic is relevant.`, description, formatEpicsForAI(epics))

		fmt.Println(" Finding the best epic match...")
		epicMatch, err := aiClient.Ask(context.Background(), epicPrompt)
		if err != nil {
			return fmt.Errorf("failed to match epic: %w", err)
		}

		epicMatch = strings.TrimSpace(epicMatch)
		if epicMatch != "NONE" {
			for _, epic := range epics {
				if epic.ID == epicMatch {
					targetEpic = &epic
					break
				}
			}
		}
	}

	// Generate detailed task description
	var epicContext string
	if targetEpic != nil {
		epicContext = targetEpic.Description
	} else {
		epicContext = "No specific epic context"
	}
	
	taskPrompt := fmt.Sprintf(`Create a detailed task description for:
"%s"

The task should be specific, actionable, and include clear success criteria.

If this epic is relevant, consider its context:
%s

Respond with a concise but detailed description that includes:
1. The specific work to be done
2. Any technical considerations
3. Clear success criteria`, description, epicContext)

	fmt.Println(" Generating detailed task description...")
	taskContent, err := aiClient.Ask(context.Background(), taskPrompt)
	if err != nil {
		return fmt.Errorf("failed to generate task content: %w", err)
	}

	// Create task file
	status, _ := cmd.Flags().GetString("status")
	taskID := utils.GenerateID("T")
	taskPath := filepath.Join("yolo", "tasks", fmt.Sprintf("%s.md", taskID))

	var epicHeader, parentLink string
	if targetEpic != nil {
		epicHeader = fmt.Sprintf("Epic: [%s] %s", targetEpic.ID, targetEpic.Title)
		parentLink = fmt.Sprintf("[%s] %s", targetEpic.ID, targetEpic.Title)
	} else {
		epicHeader = "Epic: None"
		parentLink = "None"
	}

	taskFileContent := fmt.Sprintf(`# [%s] %s

## Status: %s
Created: %s
Last Updated: %s
%s

## Description
%s

## Success Criteria
- [ ] Task implemented
- [ ] Code reviewed
- [ ] Tests added
- [ ] Documentation updated

## Relationships
<!-- YOLO-LINKS-START -->
- Parent: %s
<!-- YOLO-LINKS-END -->
`, taskID, description, status, 
   time.Now().Format("2006-01-02"), 
   time.Now().Format("2006-01-02"),
   epicHeader,
   strings.TrimSpace(taskContent),
   parentLink)

	if err := os.MkdirAll(filepath.Dir(taskPath), 0755); err != nil {
		return fmt.Errorf("failed to create tasks directory: %w", err)
	}

	if err := os.WriteFile(taskPath, []byte(taskFileContent), 0644); err != nil {
		return fmt.Errorf("failed to write task file: %w", err)
	}

	// Update epic file if we have a target epic
	if targetEpic != nil {
		if err := updateEpicRelationships(targetEpic.Path, taskID, description); err != nil {
			return fmt.Errorf("failed to update epic relationships: %w", err)
		}
	}

	fmt.Printf("\n Task %s created successfully!\n", taskID)
	if targetEpic != nil {
		fmt.Printf(" Linked to epic: [%s] %s\n", targetEpic.ID, targetEpic.Title)
	}

	return nil
}

func loadEpics() ([]Epic, error) {
	epicsDir := filepath.Join("yolo", "epics")
	files, err := filepath.Glob(filepath.Join(epicsDir, "*.md"))
	if err != nil {
		return nil, err
	}

	var epics []Epic
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		contentStr := string(content)
		idMatch := regexp.MustCompile(`# \[(E\d+)\] (.+)`).FindStringSubmatch(contentStr)
		if len(idMatch) < 3 {
			continue
		}

		descriptionMatch := regexp.MustCompile(`## Description\n((?s).+?)(?:\n#|$)`).FindStringSubmatch(contentStr)
		description := ""
		if len(descriptionMatch) >= 2 {
			description = strings.TrimSpace(descriptionMatch[1])
		}

		statusMatch := regexp.MustCompile(`## Status: (.+)`).FindStringSubmatch(contentStr)
		status := "unknown"
		if len(statusMatch) >= 2 {
			status = strings.TrimSpace(statusMatch[1])
		}

		epics = append(epics, Epic{
			ID:          idMatch[1],
			Title:       strings.TrimSpace(idMatch[2]),
			Description: description,
			Status:      status,
			Path:        file,
			Content:     contentStr,
		})
	}

	return epics, nil
}

func formatEpicsForAI(epics []Epic) string {
	var sb strings.Builder
	for _, epic := range epics {
		fmt.Fprintf(&sb, "[%s] %s\nStatus: %s\n%s\n\n", 
			epic.ID, epic.Title, epic.Status, epic.Description)
	}
	return sb.String()
}

func updateEpicRelationships(epicPath, taskID, taskTitle string) error {
	content, err := os.ReadFile(epicPath)
	if err != nil {
		return err
	}

	contentStr := string(content)
	relationSection := "## Relationships\n<!-- YOLO-LINKS-START -->\n"
	
	// Check if relationships section exists
	if !strings.Contains(contentStr, relationSection) {
		// Add it before the last section
		contentStr = strings.TrimSpace(contentStr) + "\n\n" + relationSection + "<!-- YOLO-LINKS-END -->\n"
	}

	// Update the relationships section
	re := regexp.MustCompile(`(## Relationships\n<!-- YOLO-LINKS-START -->.*?)(<!-- YOLO-LINKS-END -->)`)
	updatedContent := re.ReplaceAllStringFunc(contentStr, func(match string) string {
		// Extract existing content
		existing := re.FindStringSubmatch(match)[1]
		// Add new task if not already present
		taskLink := fmt.Sprintf("- Task: [%s] %s\n", taskID, taskTitle)
		if !strings.Contains(existing, taskLink) {
			existing += taskLink
		}
		return existing + "<!-- YOLO-LINKS-END -->"
	})

	return os.WriteFile(epicPath, []byte(updatedContent), 0644)
}