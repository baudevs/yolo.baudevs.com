package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/baudevs/yolo.baudevs.com/internal/templates"
	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
)

type ProjectStructure struct {
	Directories []string
	Files       map[string]string
}

func InitializeProject() error {
	structure := ProjectStructure{
		Directories: []string{
			"yolo",
			"yolo/epics",
			"yolo/features",
			"yolo/tasks",
			"yolo/relationships",
			"yolo/settings",
		},
		Files: map[string]string{
			"HISTORY.yml":          templates.HistoryTemplate,
			"CHANGELOG.md":         templates.ChangelogTemplate,
			"README.md":           templates.ReadmeTemplate,
			"WISHES.md":           templates.WishesTemplate,
			"STRATEGY.md":         templates.StrategyTemplate,
			"yolo/README.md":      templates.YoloReadmeTemplate,
			"LLM_INSTRUCTIONS.md": templates.LLMInstructionsTemplate,
			"yolo/settings/prompts.yml": templates.PromptsTemplate,
		},
	}

	success := color.New(color.FgGreen).SprintFunc()
	info := color.New(color.FgCyan).SprintFunc()

	fmt.Println(info("🚀 Initializing YOLO methodology..."))

	// Create directories
	for _, dir := range structure.Directories {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
		fmt.Printf("%s Created directory: %s\n", success("✓"), dir)
	}

	// Create files
	for path, template := range structure.Files {
		content := os.Expand(template, func(key string) string {
			if key == "date" {
				return time.Now().Format("2006-01-02")
			}
			return ""
		})

		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create file %s: %w", path, err)
		}
		fmt.Printf("%s Created file: %s\n", success("✓"), path)
	}

	// Check for existing git history
	versions, err := ParseGitHistory()
	if err == nil && len(versions) > 0 {
		info := color.New(color.FgCyan).SprintFunc()
		fmt.Println(info("\n📚 Found git history, converting to YOLO format..."))
		
		// Update history template with git data
		historyContent := generateHistoryFromGit(versions)
		structure.Files["HISTORY.yml"] = historyContent
		
		// Update changelog
		changelogContent := generateChangelogFromGit(versions)
		structure.Files["CHANGELOG.md"] = changelogContent
	}

	fmt.Println(success("\n✨ YOLO methodology initialized successfully!"))
	fmt.Println(info("\nNext steps:"))
	fmt.Println("1. Review README.md for project overview")
	fmt.Println("2. Check LLM_INSTRUCTIONS.md for AI developer guidelines")
	fmt.Println("3. Start creating epics with: yolo epic create")
	fmt.Println("4. Add features with: yolo feature create")
	fmt.Println("5. Create tasks with: yolo task create")

	return nil
}

func generateHistoryFromGit(versions []YoloVersion) string {
	yamlData, err := yaml.Marshal(versions)
	if err != nil {
		return templates.HistoryTemplate
	}
	return string(yamlData)
}

func generateChangelogFromGit(versions []YoloVersion) string {
	var changelog strings.Builder
	changelog.WriteString("# Changelog\n\n")

	for _, version := range versions {
		changelog.WriteString(fmt.Sprintf("## [%s] - %s\n\n", version.Version, version.Date))

		// Group changes by type
		changesByType := make(map[string][]YoloChange)
		for _, change := range version.Changes {
			changesByType[change.Type] = append(changesByType[change.Type], change)
		}

		for _, changeType := range []string{"feat", "fix", "docs", "style", "refactor", "perf", "test", "build", "ci", "chore"} {
			changes := changesByType[changeType]
			if len(changes) > 0 {
				changelog.WriteString(fmt.Sprintf("### %s\n", strings.Title(changeType)))
				for _, change := range changes {
					changelog.WriteString(fmt.Sprintf("- %s\n", change.Description))
					if change.Impact != "" {
						changelog.WriteString(fmt.Sprintf("  - Impact: %s\n", change.Impact))
					}
				}
				changelog.WriteString("\n")
			}
		}
	}

	return changelog.String()
} 