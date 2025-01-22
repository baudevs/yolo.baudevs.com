package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/baudevs/yolo.baudevs.com/internal/config"
	"github.com/baudevs/yolo.baudevs.com/internal/messages"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type MessagePrompts struct {
	Version  string                    `yaml:"version"`
	Date     string                    `yaml:"date"`
	Messages map[string]messages.Message `yaml:"messages"`
}

// InitMessagePromptsCmd returns the message prompts command
func InitMessagePromptsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "messages",
		Short: "Manage message prompts",
		Long:  "Manage YOLO's message prompts for different contexts",
	}

	cmd.AddCommand(EditPersonalityCmd())
	cmd.AddCommand(ResetPersonalityCmd())
	cmd.AddCommand(EditMethodologyCmd())
	cmd.AddCommand(ResetMethodologyCmd())

	return cmd
}

// EditPersonalityCmd returns the edit personality command
func EditPersonalityCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "edit-personality",
		Short: "Edit personality messages",
		Long:  "Edit YOLO's personality-based messages using your default editor",
		RunE: func(cmd *cobra.Command, args []string) error {
			configDir := filepath.Join(filepath.Join(os.Getenv("HOME"), ".config"), "yolo", "prompts.yml")

			// Create or load existing prompts
			prompts := &MessagePrompts{
				Version:  "1.0.0",
				Date:     "2025-01-21",
				Messages: messages.DefaultMessages,
			}

			if err := loadYAML(configDir, prompts); err != nil {
				return err
			}

			// Edit prompts
			if err := editFile(configDir, prompts); err != nil {
				return err
			}

			fmt.Println("✨ Personality messages updated!")
			return nil
		},
	}
}

// ResetPersonalityCmd returns the reset personality command
func ResetPersonalityCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reset-personality",
		Short: "Reset personality messages",
		Long:  "Reset all personality messages to their default values",
		RunE: func(cmd *cobra.Command, args []string) error {
			configDir := filepath.Join(filepath.Join(os.Getenv("HOME"), ".config"), "yolo", "prompts.yml")

			prompts := &MessagePrompts{
				Version:  "1.0.0",
				Date:     "2025-01-21",
				Messages: messages.DefaultMessages,
			}

			if err := saveYAML(configDir, prompts); err != nil {
				return err
			}

			fmt.Println("✨ Personality messages reset to defaults!")
			return nil
		},
	}
}

// EditMethodologyCmd returns the edit methodology command
func EditMethodologyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "edit-methodology",
		Short: "Edit methodology prompts",
		Long:  "Edit YOLO's methodology prompts using your default editor",
		RunE: func(cmd *cobra.Command, args []string) error {
			configDir := filepath.Join(filepath.Join(os.Getenv("HOME"), ".config"), "yolo", "methodology_prompts.yml")

			// Load or create prompts
			prompts, err := loadMethodologyPrompts()
			if err != nil {
				return err
			}

			// Edit prompts
			if err := editFile(configDir, prompts); err != nil {
				return err
			}

			fmt.Println("✨ Methodology prompts updated!")
			return nil
		},
	}
}

// ResetMethodologyCmd returns the reset methodology command
func ResetMethodologyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reset-methodology",
		Short: "Reset methodology prompts",
		Long:  "Reset all methodology prompts to their default values",
		RunE: func(cmd *cobra.Command, args []string) error {
			configDir := filepath.Join(filepath.Join(os.Getenv("HOME"), ".config"), "yolo", "methodology_prompts.yml")

			prompts, err := loadMethodologyPrompts()
			if err != nil {
				return err
			}

			if err := saveYAML(configDir, prompts); err != nil {
				return err
			}

			fmt.Println("✨ Methodology prompts reset to defaults!")
			return nil
		},
	}
}

func NewPromptCommand() *cobra.Command {
	promptCmd := &cobra.Command{
		Use:   "prompt",
		Short: "Manage prompts",
		Long:  "Manage AI prompts and personalities",
	}

	promptCmd.AddCommand(
		newPromptListCommand(),
		newPromptSetCommand(),
		newPromptPersonalityCommand(),
	)

	return promptCmd
}

func newPromptListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available prompts",
		Long:  "List all available prompts and their current values",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load client config
			clientConfig, err := config.LoadClientConfig()
			if err != nil {
				return fmt.Errorf("failed to load client config: %w", err)
			}

			fmt.Println("📝 Available Prompts")
			fmt.Println("-----------------")

			if len(clientConfig.Prompts) == 0 {
				fmt.Println("No prompts configured")
				return nil
			}

			for key, value := range clientConfig.Prompts {
				fmt.Printf("%s:\n%s\n\n", key, value)
			}

			return nil
		},
	}
}

func newPromptSetCommand() *cobra.Command {
	var key, value string

	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set a prompt",
		Long:  "Set the value for a specific prompt",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load client config
			clientConfig, err := config.LoadClientConfig()
			if err != nil {
				return fmt.Errorf("failed to load client config: %w", err)
			}

			// Initialize prompts map if needed
			if clientConfig.Prompts == nil {
				clientConfig.Prompts = make(map[string]string)
			}

			// Set prompt
			clientConfig.Prompts[key] = value

			// Save config
			if err := config.SaveClientConfig(clientConfig); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			fmt.Printf("✅ Prompt '%s' updated successfully!\n", key)
			return nil
		},
	}

	cmd.Flags().StringVarP(&key, "key", "k", "", "Prompt key (e.g., commit, error, ask)")
	cmd.Flags().StringVarP(&value, "value", "v", "", "Prompt value")
	cmd.MarkFlagRequired("key")
	cmd.MarkFlagRequired("value")

	return cmd
}

func newPromptPersonalityCommand() *cobra.Command {
	var personality string

	cmd := &cobra.Command{
		Use:   "personality",
		Short: "Set AI personality",
		Long:  "Set the AI personality for responses",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get personality level
			level := messages.GetPersonalityFromString(personality)
			if level == messages.Unknown {
				return fmt.Errorf("invalid personality: %s", personality)
			}

			// Set personality
			messages.SetPersonality(level)

			fmt.Printf("✅ AI personality set to: %s\n", personality)
			return nil
		},
	}

	cmd.Flags().StringVarP(&personality, "type", "t", "", "Personality type (nerdy, rude, unhinged)")
	cmd.MarkFlagRequired("type")

	return cmd
}

func loadYAML(file string, v interface{}) error {
	data, err := os.ReadFile(file)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if len(data) == 0 {
		return nil
	}
	return yaml.Unmarshal(data, v)
}

func saveYAML(file string, v interface{}) error {
	data, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	return os.WriteFile(file, data, 0644)
}

func editFile(file string, content interface{}) error {
	// Save current content to file
	if err := saveYAML(file, content); err != nil {
		return err
	}

	// Get editor from environment or use default
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	// Open file in editor
	cmd := exec.Command(editor, file)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func loadMethodologyPrompts() (*MessagePrompts, error) {
	return &MessagePrompts{
		Version: "1.0.0",
		Date:    "2025-01-21",
		Messages: map[string]messages.Message{
			"epic": {
				NerdyClean: `Create a new epic:
- Title
- Description
- Success Metrics
- Timeline
- Related Features`,
				MildlyRude: `Time to create an epic, boss:
- Give it a name
- What's it about?
- How do we know it's not garbage?
- When's it due?
- What other stuff needs doing?`,
				UnhingedFunny: `EPIC TIME! 🚀
- Drop that epic title like it's hot! 🔥
- Spill the tea on what this beast does! ☕️
- How do we know when we've crushed it? 💪
- When's this bad boy going live? ⏰
- What other cool stuff is coming along for the ride? 🎢`,
			},
			"feature": {
				NerdyClean: `Add a feature:
- Feature Name
- Description
- Implementation Tasks
- Dependencies`,
				MildlyRude: `New feature incoming:
- What're we calling it?
- What's it do?
- What needs doing?
- What's blocking us?`,
				UnhingedFunny: `FEATURE PARTY TIME! 🎉
- Give this bad boy a name! 🏷️
- What magic does it do? ✨
- Break it down for the peasants! 👑
- What's gonna make us cry? 😭`,
			},
			"task": {
				NerdyClean: `Create task:
- Task Name
- Implementation Details
- Acceptance Criteria
- Dependencies`,
				MildlyRude: `Task time:
- Name this thing
- How do we build it?
- When's it done?
- What's in the way?`,
				UnhingedFunny: `TASK ATTACK! 💥
- Name this beast! 🦁
- How do we make it happen? 🛠️
- When can we pop the champagne? 🍾
- What's trying to stop us? 🚧`,
			},
		},
	}, nil
}
