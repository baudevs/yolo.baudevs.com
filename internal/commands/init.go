package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type InitOptions struct {
	ProjectName           string
	ProjectPath           string
	UseGit               bool
	UseConventionalCommits bool
	FolderStructure      []string
	CustomPrompts        bool
	AIProvider           string
}

type ProjectConfig struct {
	ProjectName           string   `yaml:"project_name"`
	AIProvider           string   `yaml:"ai_provider"`
	UseConventionalCommits bool   `yaml:"use_conventional_commits"`
	CustomPrompts        bool     `yaml:"custom_prompts"`
	FolderStructure      []string `yaml:"folder_structure"`
}

// InitCmd returns the init command
func InitCmd() *cobra.Command {
	var opts InitOptions
	
	cmd := &cobra.Command{
		Use:   "init [project-name]",
		Short: "Initialize a new YOLO project",
		Long: `Initialize a new project with YOLO methodology.
This will create the necessary directory structure and configuration files.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get project name from args or current directory
			if len(args) > 0 {
				opts.ProjectName = args[0]
				opts.ProjectPath = args[0]
			} else {
				cwd, err := os.Getwd()
				if err != nil {
					return fmt.Errorf("failed to get current directory: %w", err)
				}
				opts.ProjectPath = "."
				opts.ProjectName = filepath.Base(cwd)
			}

			// Set defaults
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

			// Check if project already exists
			if existingConfig, err := detectExistingProject(opts.ProjectPath); err == nil {
				fmt.Printf("Project %s already exists. Updating...\n", opts.ProjectName)
				return updateExistingProject(opts.ProjectPath, existingConfig)
			}

			// Create new project
			return createProject(opts)
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

func updateExistingProject(path string, existingConfig *ProjectConfig) error {
	// Get config directory
	configDir := getConfigDir()

	// Only update system files, never user-customized files
	systemFiles := map[string]string{
		"yolo/settings/config.yml": filepath.Join(configDir, "settings", "config.yml"),
		"yolo/settings/ai.yml":    filepath.Join(configDir, "settings", "ai.yml"),
		"yolo/settings/git.yml":   filepath.Join(configDir, "settings", "git.yml"),
	}

	// Update each system file if it exists in the source
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
	}

	// Never update user-customizable files:
	// - README.md
	// - CHANGELOG.md
	// - HISTORY.yml
	// - STRATEGY.md
	// - WISHES.md
	// - LLM_INSTRUCTIONS.md
	// - prompts.yml (user's custom prompts)
	// - messages/personality/* (user's custom personality messages)
	// - messages/methodology/* (user's custom methodology messages)
	// - Any other user-created files

	fmt.Println("\n✨ Project updated successfully!")
	fmt.Println("✓ System files updated")
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
	// Create base directory if not current directory
	if opts.ProjectPath != "." {
		if err := os.MkdirAll(opts.ProjectPath, 0755); err != nil {
			return fmt.Errorf("failed to create project directory: %w", err)
		}
	}

	// Create selected folders
	for _, folder := range opts.FolderStructure {
		path := fmt.Sprintf("%s/yolo/%s", opts.ProjectPath, folder)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create folder %s: %w", folder, err)
		}
	}

	// Create base files
	files := map[string]string{
		"README.md": fmt.Sprintf("# %s\n\nProject created with YOLO methodology\n", opts.ProjectName),
		"CHANGELOG.md": "# Changelog\n\nAll notable changes to this project will be documented in this file.\n",
		"HISTORY.yml": "version: 1\nhistory: []\n",
		"STRATEGY.md": "# Project Strategy\n\nOutline your project strategy here.\n",
		"WISHES.md": "# Project Wishes\n\nDocument project aspirations and future plans here.\n",
		"LLM_INSTRUCTIONS.md": "# AI Interaction Guidelines\n\nGuidelines for interacting with AI assistants.\n",
	}

	for name, content := range files {
		path := fmt.Sprintf("%s/%s", opts.ProjectPath, name)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create file %s: %w", name, err)
		}
	}

	// Initialize Git if requested
	if opts.UseGit {
		if err := initGit(opts.ProjectPath); err != nil {
			return fmt.Errorf("failed to initialize git: %w", err)
		}
	}

	// Create settings directory
	settingsDir := fmt.Sprintf("%s/yolo/settings", opts.ProjectPath)
	if err := os.MkdirAll(settingsDir, 0755); err != nil {
		return fmt.Errorf("failed to create settings directory: %w", err)
	}

	// Create settings
	settings := map[string]interface{}{
		"project_name": opts.ProjectName,
		"ai_provider": opts.AIProvider,
		"use_conventional_commits": opts.UseConventionalCommits,
		"custom_prompts": opts.CustomPrompts,
		"folder_structure": opts.FolderStructure,
	}

	// Save settings to YAML file
	settingsData, err := yaml.Marshal(settings)
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	settingsPath := filepath.Join(settingsDir, "config.yml")
	if err := os.WriteFile(settingsPath, settingsData, 0644); err != nil {
		return fmt.Errorf("failed to write settings file: %w", err)
	}

	// Copy default prompts
	defaultPromptsFile := "/Users/juandavidarroyave/Developer/yolo.baudevs.com/yolo/yolo/settings/prompts.yml"
	promptsData, err := os.ReadFile(defaultPromptsFile)
	if err != nil {
		return fmt.Errorf("failed to read default prompts: %w", err)
	}

	promptsPath := filepath.Join(settingsDir, "prompts.yml")
	if err := os.WriteFile(promptsPath, promptsData, 0644); err != nil {
		return fmt.Errorf("failed to copy default prompts: %w", err)
	}

	// Create initial commit if using Git
	if opts.UseGit {
		if err := exec.Command("git", "-C", opts.ProjectPath, "add", ".").Run(); err != nil {
			return fmt.Errorf("failed to stage files: %w", err)
		}

		commitMsg := "feat: initialize project with YOLO methodology\n\n" +
			"- Create project structure\n" +
			"- Set up configuration\n" +
			"- Initialize documentation"

		if err := exec.Command("git", "-C", opts.ProjectPath, "commit", "-m", commitMsg).Run(); err != nil {
			return fmt.Errorf("failed to create initial commit: %w", err)
		}
	}

	fmt.Printf("\n✨ Project %s created successfully!\n", opts.ProjectName)
	if opts.UseGit {
		fmt.Println("✓ Git repository initialized")
	}
	fmt.Println("✓ Project structure created")
	fmt.Println("✓ Configuration saved")
	fmt.Println("\nGet started by:")
	if opts.ProjectPath != "." {
		fmt.Printf("  cd %s\n", opts.ProjectPath)
	}
	fmt.Println("  yolo status  # Check project status")
	fmt.Println("  yolo prompt  # Get AI assistance")

	return nil
}

func initGit(path string) error {
	// Initialize Git repository
	if err := exec.Command("git", "init", path).Run(); err != nil {
		return fmt.Errorf("failed to initialize git repository: %w", err)
	}

	// Create .gitignore
	gitignore := []string{
		".DS_Store",
		"*.log",
		"node_modules/",
		"dist/",
		"build/",
		".env",
		".env.*",
		"!.env.example",
	}

	gitignorePath := fmt.Sprintf("%s/.gitignore", path)
	if err := os.WriteFile(gitignorePath, []byte(strings.Join(gitignore, "\n")), 0644); err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}

	return nil
}