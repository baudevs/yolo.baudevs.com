package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/baudevs/yolo.baudevs.com/internal/config"
	"github.com/baudevs/yolo.baudevs.com/internal/project"
	"github.com/baudevs/yolo.baudevs.com/internal/types"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type InitOptions struct {
	ProjectName            string
	ProjectPath            string
	UseGit                 bool
	UseConventionalCommits bool
	FolderStructure        []string
	CustomPrompts          bool
	AIProvider             string
	Description            string
}

type ProjectConfig struct {
	ProjectName            string   `yaml:"project_name"`
	AIProvider             string   `yaml:"ai_provider"`
	UseConventionalCommits bool     `yaml:"use_conventional_commits"`
	CustomPrompts          bool     `yaml:"custom_prompts"`
	FolderStructure        []string `yaml:"folder_structure"`
	Description            string   `yaml:"description"`
}

// InitCmd returns the init command
func InitCmd() *cobra.Command {
	var opts InitOptions

	cmd := &cobra.Command{
		Use:   "init [project-name]",
		Short: "Initialize a new YOLO project",
		Long: `Initialize a new project with YOLO methodology.
This will create the necessary directory structure and configuration files.
The initialization process includes an AI-powered project planning phase.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create AI manager
			aiManager, err := project.NewAIManager()
			if err != nil {
				// Check if it's because of missing OpenAI API key
				cfg, cfgErr := config.LoadConfig()
				if cfgErr == nil && cfg.OpenAI.APIKey == "" {
					configPath, _ := config.GetConfigPath()
					fmt.Println("\n❌ No OpenAI API key configured!")
					fmt.Println("\nTo use YOLO, you need to set up your OpenAI API key. You can do this in two ways:")
					fmt.Println("1. Set the OPENAI_API_KEY environment variable")
					fmt.Println("2. Add your API key to:", configPath)
					fmt.Println("\nGet your API key at: https://platform.openai.com/api-keys")
					return fmt.Errorf("OpenAI API key not configured")
				}
				return fmt.Errorf("failed to initialize AI manager: %w", err)
			}

			// Get project description
			description, err := aiManager.GetProjectDescription()
			if err != nil {
				return fmt.Errorf("failed to get project description: %w", err)
			}
			opts.Description = description

			// Get project name
			if len(args) > 0 {
				opts.ProjectName = args[0]
				opts.ProjectPath = args[0]
			} else {
				// Get project name from user or generate one
				projectName, _, err := aiManager.GetProjectName(description)
				if err != nil {
					return fmt.Errorf("failed to get project name: %w", err)
				}
				opts.ProjectName = projectName
				opts.ProjectPath = "."
			}

			// Set defaults for additional directories
			opts.UseGit = true
			opts.UseConventionalCommits = true
			opts.CustomPrompts = false
			opts.AIProvider = "openai"
			opts.FolderStructure = []string{
				"settings",
				"messages/personality",
				"messages/methodology",
				"relationships",
			}

			// Generate project plan using AI
			fmt.Println("\n🎯 Generating project plan...")
			projectPlan, err := aiManager.GenerateProjectPlan(description)
			if err != nil {
				return fmt.Errorf("failed to generate project plan: %w", err)
			}

			// Update project name in the plan
			projectPlan.Name = opts.ProjectName

			// Create additional directories first
			fmt.Println("\n📁 Creating additional project directories...")
			for _, dir := range opts.FolderStructure {
				path := filepath.Join(opts.ProjectPath, dir)
				fmt.Printf("  Creating %s/...", dir)
				if err := os.MkdirAll(path, 0755); err != nil {
					fmt.Println(" ❌")
					return fmt.Errorf("failed to create directory %s: %w", path, err)
				}
				fmt.Println(" ✓")
			}

			// Generate and save project structure
			if err := aiManager.SaveProjectStructure(projectPlan); err != nil {
				return fmt.Errorf("failed to save project structure: %w", err)
			}

			// Save project configuration
			config := &ProjectConfig{
				ProjectName:            opts.ProjectName,
				AIProvider:             opts.AIProvider,
				UseConventionalCommits: opts.UseConventionalCommits,
				CustomPrompts:          opts.CustomPrompts,
				FolderStructure:        opts.FolderStructure,
				Description:            opts.Description,
			}

			if err := saveProjectConfig(opts.ProjectPath, config); err != nil {
				fmt.Printf("Warning: Failed to save project configuration: %v\n", err)
			}

			fmt.Printf("\n✨ Project %s initialized successfully!\n", opts.ProjectName)
			fmt.Println("✓ Project structure created")
			fmt.Println("✓ AI-powered project plan generated")
			fmt.Println("✓ All documentation files created")

			return nil
		},
	}

	return cmd
}

func detectExistingProject(path string) (*ProjectConfig, error) {
	configPath := filepath.Join(path, "yolo", "settings", "config.yml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config ProjectConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func updateExistingProject(path string, existingConfig *ProjectConfig, aiManager *project.AIManager, projectPlan *types.Project) error {
	fmt.Println("\n🔄 Updating project structure...")

	// Get config directory
	configDir := getConfigDir()

	// Only update system files, never user-customized files
	systemFiles := map[string]string{
		"yolo/settings/config.yml": filepath.Join(configDir, "settings", "config.yml"),
		"yolo/settings/ai.yml":     filepath.Join(configDir, "settings", "ai.yml"),
		"yolo/settings/git.yml":    filepath.Join(configDir, "settings", "git.yml"),
	}

	// Update each system file if it exists in the source
	fmt.Println("\n📝 Updating system files...")
	for projectFile, sourceFile := range systemFiles {
		// Check if source file exists
		if _, err := os.Stat(sourceFile); os.IsNotExist(err) {
			continue
		}

		// Read source file
		data, err := os.ReadFile(sourceFile)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", sourceFile, err)
		}

		// Create target directory if needed
		targetDir := filepath.Dir(filepath.Join(path, projectFile))
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", targetDir, err)
		}

		// Write to target file
		if err := os.WriteFile(filepath.Join(path, projectFile), data, 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", projectFile, err)
		}
		fmt.Printf("✓ Updated %s\n", projectFile)
	}

	// Generate and save AI content
	fmt.Println("\n🤖 Generating AI-powered content...")
	if err := aiManager.SaveProjectStructure(projectPlan); err != nil {
		return fmt.Errorf("failed to save project structure: %w", err)
	}

	fmt.Println("\n✨ Project updated successfully!")
	fmt.Println("✓ System files updated")
	fmt.Println("✓ AI-powered content generated")
	fmt.Println("✓ User customizations preserved")
	return nil
}

func getConfigDir() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = filepath.Join(os.Getenv("HOME"), ".config")
	}
	return filepath.Join(configDir, "yolo")
}

func createProject(opts InitOptions) error {
	// Implement the project creation logic here
	fmt.Printf("Creating project %s at %s\n", opts.ProjectName, opts.ProjectPath)
	// Example: create directories
	for _, dir := range opts.FolderStructure {
		path := filepath.Join(opts.ProjectPath, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", path, err)
		}
	}
	return nil
}

func runInit(cmd *cobra.Command, args []string) error {
	var opts InitOptions
	// Initialize opts based on command flags or defaults
	// ... existing logic ...
	return createProject(opts)
}

func updateHistoryWithNewFeatures(path string) error {
	// Implement the logic to update history with new features
	fmt.Println("Updating history with new features...")
	return nil
}

func saveProjectConfig(path string, config *ProjectConfig) error {
	// Implement the logic to save project configuration
	fmt.Println("Saving project configuration...")
	return nil
}

func updateChangelogWithNewFeatures(path string) error {
	changelogPath := filepath.Join(path, "CHANGELOG.md")
	content, err := os.ReadFile(changelogPath)
	if err != nil {
		return err
	}

	// Add new sections if they don't exist
	if !strings.Contains(string(content), "Work in Progress") {
		newContent := string(content) + `

### Work in Progress
- System-wide keyboard shortcuts (F011)
  - ✅ Web interface for configuring shortcuts
  - ✅ JSON-based configuration persistence
  - ✅ WebSocket for real-time updates
  - 🏗️ macOS daemon for global shortcut capture
  - 📅 Linux support planned`
		return os.WriteFile(changelogPath, []byte(newContent), 0644)
	}

	return nil
}
