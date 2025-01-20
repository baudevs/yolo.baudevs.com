package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/baudevs/yolo.baudevs.com/internal/messages"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type MessagePrompts struct {
	Messages map[string]messages.Message `yaml:"messages"`
}

type MethodologyPrompts struct {
	StandardDocumentation string `yaml:"standard_documentation"`
	UpdateChangelog      string `yaml:"update_changelog"`
	UpdateReadme         string `yaml:"update_readme"`
	EpicDocumentation    string `yaml:"epic_documentation"`
	FeatureDocumentation string `yaml:"feature_documentation"`
	TaskDocumentation    string `yaml:"task_documentation"`
	UpdateHistory        string `yaml:"update_history"`
}

// InitMessagePromptsCmd initializes the message prompts command
func InitMessagePromptsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "messages",
		Short: "Manage YOLO's messages and prompts",
		Long: `Manage both YOLO's personality messages and methodology prompts.
		
Two types of messages:
1. Personality Messages - Used for YOLO's interactive responses
2. Methodology Prompts - Used for documentation and workflow templates`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Choose what to manage:")
			fmt.Println("  personality")
			fmt.Println("    edit     Edit personality messages")
			fmt.Println("    reset    Reset personality messages to defaults")
			fmt.Println("  methodology")
			fmt.Println("    edit     Edit methodology prompts")
			fmt.Println("    reset    Reset methodology prompts to defaults")
		},
	}

	personalityCmd := &cobra.Command{
		Use:   "personality",
		Short: "Manage YOLO's personality messages",
		Long:  "Edit or reset YOLO's personality-based messages and responses",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Use one of the subcommands:")
			fmt.Println("  edit    Edit personality messages")
			fmt.Println("  reset   Reset personality messages to defaults")
		},
	}

	methodologyCmd := &cobra.Command{
		Use:   "methodology",
		Short: "Manage YOLO's methodology prompts",
		Long:  "Edit or reset YOLO's documentation and workflow prompts",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Use one of the subcommands:")
			fmt.Println("  edit    Edit methodology prompts")
			fmt.Println("  reset   Reset methodology prompts to defaults")
		},
	}

	personalityCmd.AddCommand(
		editPersonalityCmd(),
		resetPersonalityCmd(),
	)

	methodologyCmd.AddCommand(
		editMethodologyCmd(),
		resetMethodologyCmd(),
	)

	cmd.AddCommand(personalityCmd, methodologyCmd)
	return cmd
}

func editPersonalityCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "edit",
		Short: "Edit personality messages",
		Long:  "Edit YOLO's personality-based messages using your default editor",
		RunE: func(cmd *cobra.Command, args []string) error {
			configDir := getConfigDir()
			promptsFile := filepath.Join(configDir, "prompts.yml")

			// Create or load existing prompts
			prompts, err := loadMessagePrompts(promptsFile)
			if err != nil {
				prompts = &MessagePrompts{
					Messages: messages.DefaultMessages,
				}
			}

			return editFile(promptsFile, prompts)
		},
	}
}

func resetPersonalityCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reset",
		Short: "Reset personality messages",
		Long:  "Reset all personality messages to their default values",
		RunE: func(cmd *cobra.Command, args []string) error {
			configDir := getConfigDir()
			promptsFile := filepath.Join(configDir, "prompts.yml")

			prompts := &MessagePrompts{
				Messages: messages.DefaultMessages,
			}

			if err := saveYAML(promptsFile, prompts); err != nil {
				return err
			}

			fmt.Println("🔄 Personality messages have been reset to defaults!")
			return nil
		},
	}
}

func editMethodologyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "edit",
		Short: "Edit methodology prompts",
		Long:  "Edit YOLO's methodology prompts using your default editor",
		RunE: func(cmd *cobra.Command, args []string) error {
			configDir := getConfigDir()
			promptsFile := filepath.Join(configDir, "methodology_prompts.yml")

			// Load or create prompts
			prompts, err := loadMethodologyPrompts()
			if err != nil {
				return fmt.Errorf("failed to load methodology prompts: %w", err)
			}

			return editFile(promptsFile, prompts)
		},
	}
}

func resetMethodologyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reset",
		Short: "Reset methodology prompts",
		Long:  "Reset all methodology prompts to their default values",
		RunE: func(cmd *cobra.Command, args []string) error {
			configDir := getConfigDir()
			promptsFile := filepath.Join(configDir, "methodology_prompts.yml")

			// Create default prompts
			prompts, err := loadMethodologyPrompts()
			if err != nil {
				return fmt.Errorf("failed to load default methodology prompts: %w", err)
			}

			if err := saveYAML(promptsFile, prompts); err != nil {
				return err
			}

			fmt.Println("🔄 Methodology prompts have been reset to defaults!")
			return nil
		},
	}
}

func editFile(file string, content interface{}) error {
	// Convert to YAML
	data, err := yaml.Marshal(content)
	if err != nil {
		return fmt.Errorf("error marshaling content: %v", err)
	}

	// Write to file
	if err := os.WriteFile(file, data, 0644); err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}

	// Open in default editor
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	editorCmd := exec.Command(editor, file)
	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdout
	editorCmd.Stderr = os.Stderr

	if err := editorCmd.Run(); err != nil {
		return fmt.Errorf("error opening editor: %v", err)
	}

	fmt.Println("✨ File updated successfully!")
	return nil
}

func saveYAML(file string, content interface{}) error {
	data, err := yaml.Marshal(content)
	if err != nil {
		return fmt.Errorf("error marshaling content: %v", err)
	}

	if err := os.WriteFile(file, data, 0644); err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}

	return nil
}

func getConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		os.Exit(1)
	}
	configDir := filepath.Join(home, ".config", "yolo")
	os.MkdirAll(configDir, 0755)
	return configDir
}

func loadMessagePrompts(file string) (*MessagePrompts, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("prompts file not found")
		}
		return nil, fmt.Errorf("error reading prompts: %v", err)
	}

	var prompts MessagePrompts
	if err := yaml.Unmarshal(data, &prompts); err != nil {
		return nil, fmt.Errorf("error parsing prompts: %v", err)
	}

	return &prompts, nil
}

func loadMethodologyPrompts() (*MethodologyPrompts, error) {
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
		defaultPrompts := MethodologyPrompts{
			StandardDocumentation: `# Write documentation that follows YOLO methodology:

## Core Documentation Guidelines
- Clear and concise explanations
- Code examples where relevant
- Usage instructions
- Common pitfalls and solutions

## File Organization
- Store all documentation in the @yolo directory
- Each epic, task, and feature gets its own markdown file
- Link items using @ tags in titles and subtitles
  Example: "# Task: Add OAuth [@epic/auth] [@feature/login]"

## Required Files
### @README.md
- Add new sections at the bottom
- Mark status with tags: [Implemented], [Deleted], [Deprecated], [In Progress], [WIP], [Frozen]
- Never delete content, only update status

### @CHANGELOG.md
- Add entries in reverse chronological order
- Include status changes and @ references
- Never remove entries

### @HISTORY.yaml
- Maintain chronological order of all actions
- Document both done and undone items
- Include timestamps and @ references

### @WISHES.md
- Track feature requests and improvements
- Update status of existing wishes
- Cross-reference with @ tags`,
			UpdateChangelog: `# Update changelog following YOLO methodology:

## Entry Format
- Add entries at the top (reverse chronological)
- Group by type: [Added], [Changed], [Fixed], [Deprecated], [Removed]
- Include @ references to related items
- Never delete entries, mark as [Reverted] if needed

## Required Information
- Version numbers (if applicable)
- Timestamp of change
- Author/contributor
- Links to related @yolo files
- Status tags for referenced items`,
			UpdateReadme: `# Update README following YOLO methodology:

## Content Guidelines
- Add new sections at the bottom
- Never delete existing sections
- Update status tags on existing content
- Cross-reference with @ tags

## Required Sections
- Project overview
- Installation guide
- Configuration options
- Usage examples with code
- Links to @yolo documentation
- Contributing guidelines
- Current project status`,
			EpicDocumentation: `# Document epics following YOLO methodology:

## File Location
- Create in @yolo/epics/{epic-name}.md
- Link related features with @ tags
- Reference parent epics if applicable

## Required Sections
- Overview and business value
- Success criteria and timeline
- Dependencies and risks
- Current status with timestamp
- List of child features/tasks
- Implementation history`,
			FeatureDocumentation: `# Document features following YOLO methodology:

## File Location
- Create in @yolo/features/{feature-name}.md
- Reference parent epic with @ tag
- Link related tasks with @ tags

## Required Sections
- Technical requirements
- Implementation details
- Testing strategy
- Current status and history
- Related components
- Usage examples`,
			TaskDocumentation: `# Document tasks following YOLO methodology:

## File Location
- Create in @yolo/tasks/{task-name}.md
- Reference parent feature with @ tag
- Link related tasks with @ tags

## Required Sections
- Clear objectives
- Implementation steps
- Dependencies
- Time estimates
- Current status
- Testing criteria`,
			UpdateHistory: `# Update history following YOLO methodology:

## Entry Format in @HISTORY.yaml
- Chronological order (oldest first)
- Timestamp for each entry
- Action type: [Added], [Changed], [Removed], [Reverted]
- @ references to affected items
- Author/contributor
- Status change details
- Never delete entries`,
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

	var prompts MethodologyPrompts
	if err := yaml.Unmarshal(data, &prompts); err != nil {
		return nil, fmt.Errorf("failed to parse prompts file: %w", err)
	}

	return &prompts, nil
}
