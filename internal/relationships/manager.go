package relationships

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/baudevs/yolo.baudevs.com/internal/ai"
)

type WorkItemType string

const (
	Epic    WorkItemType = "Epic"
	Feature WorkItemType = "Feature"
	Task    WorkItemType = "Task"
)

type WorkItem struct {
	Type        WorkItemType
	ID          string
	Title       string
	Description string
	Status      string
	Path        string
	Content     string
}

type RelationshipManager struct {
	aiClient *ai.Client
}

func NewManager(aiClient *ai.Client) *RelationshipManager {
	return &RelationshipManager{
		aiClient: aiClient,
	}
}

// FindOrCreateParent finds or suggests creating a parent work item
func (m *RelationshipManager) FindOrCreateParent(ctx context.Context, itemType WorkItemType, description string, existingItems []WorkItem) (*WorkItem, bool, error) {
	var parentType WorkItemType
	switch itemType {
	case Task:
		parentType = Feature
	case Feature:
		parentType = Epic
	default:
		return nil, false, nil // Epics don't have parents
	}

	prompt := fmt.Sprintf(`Given this %s description:
"%s"

Analyze these existing %ss and determine if this %s should:
1. Be linked to an existing %s (respond with ID)
2. Need a new %s to be created (respond with "NEW")

Existing %ss:
%s

Consider:
- Scope and goals alignment
- Natural relationships
- Project structure

Respond with just the %s ID (e.g., "E001") or "NEW" if a new %s should be created.`,
		itemType, description, parentType, itemType, parentType, parentType, parentType,
		formatItemsForAI(filterByType(existingItems, parentType)), parentType, parentType)

	response, err := m.aiClient.Ask(ctx, prompt)
	if err != nil {
		return nil, false, err
	}

	response = strings.TrimSpace(response)
	if response == "NEW" {
		return nil, true, nil
	}

	// Find existing parent
	for _, item := range existingItems {
		if item.Type == parentType && item.ID == response {
			return &item, false, nil
		}
	}

	return nil, false, nil
}

// SuggestChildren suggests child items that should be created
func (m *RelationshipManager) SuggestChildren(ctx context.Context, itemType WorkItemType, description string) ([]string, error) {
	var childType WorkItemType
	switch itemType {
	case Epic:
		childType = Feature
	case Feature:
		childType = Task
	default:
		return nil, nil // Tasks don't have children
	}

	prompt := fmt.Sprintf(`Given this %s description:
"%s"

Suggest 3-5 %ss that would be needed to implement this %s. Each %s should be:
- Specific and actionable
- Clear in scope
- Contribute directly to the %s's goals

Respond with one %s title per line, no numbers or bullets.`,
		itemType, description, childType, itemType, childType, itemType, childType)

	response, err := m.aiClient.Ask(ctx, prompt)
	if err != nil {
		return nil, err
	}

	// Split response into lines and clean up
	children := []string{}
	for _, line := range strings.Split(response, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			children = append(children, line)
		}
	}

	return children, nil
}

// UpdateRelationships updates the relationships section in a markdown file
func (m *RelationshipManager) UpdateRelationships(filePath string, relationships map[WorkItemType][]WorkItem) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	contentStr := string(content)
	relationSection := "## Relationships\n<!-- YOLO-LINKS-START -->\n"

	// Check if relationships section exists
	if !strings.Contains(contentStr, relationSection) {
		contentStr = strings.TrimSpace(contentStr) + "\n\n" + relationSection + "<!-- YOLO-LINKS-END -->\n"
	}

	// Build relationships content
	var sb strings.Builder
	sb.WriteString(relationSection)

	// Add parents first
	if epics := relationships[Epic]; len(epics) > 0 {
		for _, epic := range epics {
			sb.WriteString(fmt.Sprintf("- Parent Epic: [%s] %s\n", epic.ID, epic.Title))
		}
	}
	if features := relationships[Feature]; len(features) > 0 {
		for _, feature := range features {
			sb.WriteString(fmt.Sprintf("- Parent Feature: [%s] %s\n", feature.ID, feature.Title))
		}
	}

	// Add children
	if features := relationships[Feature]; len(features) > 0 {
		sb.WriteString("\n- Features:\n")
		for _, feature := range features {
			sb.WriteString(fmt.Sprintf("  - [%s] %s\n", feature.ID, feature.Title))
		}
	}
	if tasks := relationships[Task]; len(tasks) > 0 {
		sb.WriteString("\n- Tasks:\n")
		for _, task := range tasks {
			sb.WriteString(fmt.Sprintf("  - [%s] %s\n", task.ID, task.Title))
		}
	}

	sb.WriteString("<!-- YOLO-LINKS-END -->")

	// Update the relationships section
	re := regexp.MustCompile(`(## Relationships\n<!-- YOLO-LINKS-START -->.*?)(<!-- YOLO-LINKS-END -->)`)
	updatedContent := re.ReplaceAllString(contentStr, sb.String())

	return os.WriteFile(filePath, []byte(updatedContent), 0644)
}

// LoadWorkItems loads all work items of the given types
func (m *RelationshipManager) LoadWorkItems(types ...WorkItemType) ([]WorkItem, error) {
	var items []WorkItem

	for _, itemType := range types {
		var dir string
		var prefix string
		switch itemType {
		case Epic:
			dir = "epics"
			prefix = "E"
		case Feature:
			dir = "features"
			prefix = "F"
		case Task:
			dir = "tasks"
			prefix = "T"
		}

		files, err := filepath.Glob(filepath.Join("yolo", dir, "*.md"))
		if err != nil {
			continue
		}

		for _, file := range files {
			content, err := os.ReadFile(file)
			if err != nil {
				continue
			}

			contentStr := string(content)
			idMatch := regexp.MustCompile(`# \[(`+prefix+`\d+)\] (.+)`).FindStringSubmatch(contentStr)
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

			items = append(items, WorkItem{
				Type:        itemType,
				ID:          idMatch[1],
				Title:       strings.TrimSpace(idMatch[2]),
				Description: description,
				Status:      status,
				Path:        file,
				Content:     contentStr,
			})
		}
	}

	return items, nil
}

func filterByType(items []WorkItem, itemType WorkItemType) []WorkItem {
	var filtered []WorkItem
	for _, item := range items {
		if item.Type == itemType {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func formatItemsForAI(items []WorkItem) string {
	var sb strings.Builder
	for _, item := range items {
		fmt.Fprintf(&sb, "[%s] %s\nStatus: %s\n%s\n\n",
			item.ID, item.Title, item.Status, item.Description)
	}
	return sb.String()
}
