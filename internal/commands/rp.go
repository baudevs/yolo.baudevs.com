package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Prompts struct {
	StandardDocumentation string `yaml:"standard_documentation"`
	UpdateChangelog      string `yaml:"update_changelog"`
	UpdateReadme         string `yaml:"update_readme"`
	EpicDocumentation    string `yaml:"epic_documentation"`
	FeatureDocumentation string `yaml:"feature_documentation"`
	TaskDocumentation    string `yaml:"task_documentation"`
	UpdateHistory        string `yaml:"update_history"`
}

func PromptCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prompt",
		Short: "Work with YOLO methodology prompts",
		Long: `Access and copy various YOLO methodology prompts to your clipboard.
These prompts are designed to help LLMs understand and follow YOLO methodology.`,
	}

	// Add subcommands
	cmd.AddCommand(
		promptCopyCmd("standard", "Copy the standard documentation prompt", func(p Prompts) string { return p.StandardDocumentation }),
		promptCopyCmd("changelog", "Copy the changelog update prompt", func(p Prompts) string { return p.UpdateChangelog }),
		promptCopyCmd("readme", "Copy the readme update prompt", func(p Prompts) string { return p.UpdateReadme }),
		promptCopyCmd("epic", "Copy the epic documentation prompt", func(p Prompts) string { return p.EpicDocumentation }),
		promptCopyCmd("feature", "Copy the feature documentation prompt", func(p Prompts) string { return p.FeatureDocumentation }),
		promptCopyCmd("task", "Copy the task documentation prompt", func(p Prompts) string { return p.TaskDocumentation }),
		promptCopyCmd("history", "Copy the history update prompt", func(p Prompts) string { return p.UpdateHistory }),
	)

	return cmd
}

func promptCopyCmd(use string, short string, promptSelector func(Prompts) string) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		Run: func(cmd *cobra.Command, args []string) {
			prompts, err := loadPrompts()
			if err != nil {
				fmt.Printf("Error reading prompts: %v\n", err)
				return
			}

			prompt := promptSelector(*prompts)
			if err := copyToClipboard(prompt); err != nil {
				fmt.Printf("Error copying to clipboard: %v\n", err)
				return
			}

			success := color.New(color.FgGreen).SprintFunc()
			fmt.Printf("%s %s prompt copied to clipboard!\n", success("✓"), strings.Title(use))
		},
	}
}

func loadPrompts() (*Prompts, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(home, ".config", "yolo")
	promptsFile := filepath.Join(configDir, "methodology_prompts.yml")

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	// If file doesn't exist, create it with default prompts
	if _, err := os.Stat(promptsFile); os.IsNotExist(err) {
		defaultPrompts := Prompts{
			StandardDocumentation: `# Standard Documentation Prompt
Write documentation that follows YOLO methodology:
- Clear and concise explanations
- Code examples where relevant
- Usage instructions
- Common pitfalls and solutions`,
			UpdateChangelog: `# Changelog Update Prompt
Update the changelog following YOLO methodology:
- Group changes by type (Added, Changed, Fixed)
- Include version numbers
- Add date stamps
- Reference relevant issues/PRs`,
			UpdateReadme: `# README Update Prompt
Update the README following YOLO methodology:
- Project overview and purpose
- Installation instructions
- Usage examples
- Configuration options
- Contributing guidelines`,
			EpicDocumentation: `# Epic Documentation Prompt
Document epics following YOLO methodology:
- High-level overview
- Business value
- Success criteria
- Dependencies and risks
- Timeline and milestones`,
			FeatureDocumentation: `# Feature Documentation Prompt
Document features following YOLO methodology:
- Feature description
- Technical requirements
- Implementation details
- Testing strategy
- Release notes`,
			TaskDocumentation: `# Task Documentation Prompt
Document tasks following YOLO methodology:
- Clear objectives
- Implementation steps
- Acceptance criteria
- Dependencies
- Time estimates`,
			UpdateHistory: `# History Update Prompt
Update history following YOLO methodology:
- Chronological order
- Key milestones
- Major decisions
- Lessons learned`,
		}

		data, err := yaml.Marshal(defaultPrompts)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal default prompts: %w", err)
		}

		if err := os.WriteFile(promptsFile, data, 0644); err != nil {
			return nil, fmt.Errorf("failed to write default prompts: %w", err)
		}

		return &defaultPrompts, nil
	}

	// Read existing prompts
	data, err := os.ReadFile(promptsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read prompts file: %w", err)
	}

	var prompts Prompts
	if err := yaml.Unmarshal(data, &prompts); err != nil {
		return nil, fmt.Errorf("failed to parse prompts file: %w", err)
	}

	return &prompts, nil
}

func copyToClipboard(text string) error {
	var cmd *exec.Cmd
	
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "linux":
		cmd = exec.Command("xclip", "-selection", "clipboard")
	case "windows":
		cmd = exec.Command("clip")
	default:
		return fmt.Errorf("unsupported platform")
	}

	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
} 