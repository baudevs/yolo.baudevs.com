package commands

import (
	"fmt"
	"os"
	"os/exec"
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
	data, err := os.ReadFile("yolo/settings/prompts.yml")
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