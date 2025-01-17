package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/baudevs/yolo-cli/internal/core"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os/exec"
	"path/filepath"
)

var (
	docStyle = lipgloss.NewStyle().Margin(1, 2)
	
	titleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF75B7")).
		Bold(true).
		Padding(0, 1)
		
	selectedItemStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#A550DF")).
		Bold(true)
)

type InitOptions struct {
	ProjectName        string
	ProjectPath        string
	UseGit            bool
	UseConventionalCommits bool
	FolderStructure   []string
	CustomPrompts     bool
	AIProvider        string
}

type ProjectConfig struct {
	ProjectName           string   `yaml:"project_name"`
	AIProvider           string   `yaml:"ai_provider"`
	UseConventionalCommits bool   `yaml:"use_conventional_commits"`
	CustomPrompts        bool     `yaml:"custom_prompts"`
	FolderStructure      []string `yaml:"folder_structure"`
}

type step int

const (
	stepDetectProject step = iota
	stepReinitChoice
	stepConfirmDelete
	stepReconfigureChoice
	stepProjectName
	stepProjectPath
	stepGitChoice
	stepFolderStructure
	stepAIProvider
	stepCustomPrompts
	stepConfirm
	stepDone
)

type model struct {
	options InitOptions
	currentStep step
	textInput textinput.Model
	folderList list.Model
	aiList list.Model
	configList list.Model
	existingConfig *ProjectConfig
	err error
}

type item struct {
	title string
	desc  string
	selected bool
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func detectExistingProject(path string) (*ProjectConfig, error) {
	configPath := filepath.Join(path, "yolo", "settings", "config.yml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var config ProjectConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("invalid config file: %w", err)
	}

	return &config, nil
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "my-awesome-project"
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 40

	folderItems := []list.Item{
		item{title: "epics", desc: "Track large-scale project initiatives", selected: true},
		item{title: "features", desc: "Document and track feature development", selected: true},
		item{title: "tasks", desc: "Create and monitor task progress", selected: true},
		item{title: "relationships", desc: "Track connections between items", selected: true},
		item{title: "settings", desc: "Project configuration and settings", selected: true},
		item{title: "docs", desc: "Additional documentation", selected: false},
		item{title: "assets", desc: "Project assets and resources", selected: false},
	}

	aiItems := []list.Item{
		item{title: "OpenAI", desc: "Use OpenAI's GPT models"},
		item{title: "Anthropic", desc: "Use Anthropic's Claude models"},
		item{title: "Custom", desc: "Configure your own AI provider"},
	}

	configItems := []list.Item{
		item{title: "Project Name", desc: "Change the project name"},
		item{title: "AI Provider", desc: "Change the AI provider"},
		item{title: "Folder Structure", desc: "Modify project folders"},
		item{title: "Git Settings", desc: "Update Git configuration"},
		item{title: "Custom Prompts", desc: "Enable/disable custom prompts"},
	}

	folderList := list.New(folderItems, list.NewDefaultDelegate(), 0, 0)
	folderList.Title = "Select Folder Structure"
	folderList.SetShowHelp(false)

	aiList := list.New(aiItems, list.NewDefaultDelegate(), 0, 0)
	aiList.Title = "Select AI Provider"
	aiList.SetShowHelp(false)

	configList := list.New(configItems, list.NewDefaultDelegate(), 0, 0)
	configList.Title = "What would you like to reconfigure?"
	configList.SetShowHelp(false)

	return model{
		textInput: ti,
		folderList: folderList,
		aiList: aiList,
		configList: configList,
		options: InitOptions{
			UseGit: true,
			UseConventionalCommits: true,
			FolderStructure: []string{"epics", "features", "tasks", "relationships", "settings"},
		},
		currentStep: stepDetectProject,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter":
			switch m.currentStep {
			case stepDetectProject:
				// Project detection happens automatically
				if m.existingConfig != nil {
					m.currentStep = stepReinitChoice
				} else {
					m.currentStep = stepProjectName
				}
			case stepReinitChoice:
				selectedItem := m.configList.SelectedItem().(item)
				if selectedItem.title == "Reinitialize Project" {
					m.currentStep = stepConfirmDelete
					m.textInput.Placeholder = "Type project name to confirm"
					m.textInput.SetValue("")
					m.textInput.Focus()
				} else {
					m.currentStep = stepReconfigureChoice
				}
			case stepConfirmDelete:
				if m.textInput.Value() == m.existingConfig.ProjectName {
					m.currentStep = stepProjectName
					m.textInput.Placeholder = "my-awesome-project"
					m.textInput.SetValue("")
				}
			case stepReconfigureChoice:
				selectedItem := m.configList.SelectedItem().(item)
				switch selectedItem.title {
				case "Project Name":
					m.currentStep = stepProjectName
					m.textInput.SetValue(m.existingConfig.ProjectName)
				case "AI Provider":
					m.currentStep = stepAIProvider
				case "Folder Structure":
					m.currentStep = stepFolderStructure
					// Update selected folders based on existing config
					items := make([]list.Item, 0)
					for _, listItem := range m.folderList.Items() {
						i := listItem.(item)
						i.selected = false
						for _, folder := range m.existingConfig.FolderStructure {
							if i.title == folder {
								i.selected = true
								break
							}
						}
						items = append(items, i)
					}
					m.folderList.SetItems(items)
				case "Git Settings":
					m.currentStep = stepGitChoice
				case "Custom Prompts":
					m.currentStep = stepCustomPrompts
				}
			case stepProjectName:
				if m.textInput.Value() != "" {
					m.options.ProjectName = m.textInput.Value()
					m.textInput.Placeholder = "."
					m.textInput.SetValue("")
					m.currentStep = stepProjectPath
				}
			case stepProjectPath:
				path := m.textInput.Value()
				if path == "" {
					path = "."
				}
				m.options.ProjectPath = path
				m.currentStep = stepGitChoice
			case stepGitChoice:
				m.currentStep = stepFolderStructure
			case stepFolderStructure:
				m.currentStep = stepAIProvider
			case stepAIProvider:
				selectedItem := m.aiList.SelectedItem().(item)
				m.options.AIProvider = selectedItem.title
				m.currentStep = stepCustomPrompts
			case stepCustomPrompts:
				m.currentStep = stepConfirm
			case stepConfirm:
				m.currentStep = stepDone
				return m, tea.Quit
			}

		case "y", "Y":
			if m.currentStep == stepGitChoice {
				m.options.UseGit = true
				m.currentStep = stepFolderStructure
			} else if m.currentStep == stepCustomPrompts {
				m.options.CustomPrompts = true
				m.currentStep = stepConfirm
			}

		case "n", "N":
			if m.currentStep == stepGitChoice {
				m.options.UseGit = false
				m.currentStep = stepFolderStructure
			} else if m.currentStep == stepCustomPrompts {
				m.options.CustomPrompts = false
				m.currentStep = stepConfirm
			}

		case "space":
			if m.currentStep == stepFolderStructure {
				selectedItem := m.folderList.SelectedItem().(item)
				selectedItem.selected = !selectedItem.selected
				items := make([]list.Item, 0)
				for _, listItem := range m.folderList.Items() {
					i := listItem.(item)
					if i.title == selectedItem.title {
						items = append(items, selectedItem)
					} else {
						items = append(items, i)
					}
				}
				m.folderList.SetItems(items)
			}
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.folderList.SetSize(msg.Width-h, msg.Height-v)
		m.aiList.SetSize(msg.Width-h, msg.Height-v)
		m.configList.SetSize(msg.Width-h, msg.Height-v)
	}

	if m.currentStep == stepProjectName || m.currentStep == stepProjectPath || m.currentStep == stepConfirmDelete {
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	} else if m.currentStep == stepFolderStructure {
		m.folderList, cmd = m.folderList.Update(msg)
		return m, cmd
	} else if m.currentStep == stepAIProvider {
		m.aiList, cmd = m.aiList.Update(msg)
		return m, cmd
	} else if m.currentStep == stepReinitChoice || m.currentStep == stepReconfigureChoice {
		m.configList, cmd = m.configList.Update(msg)
		return m, cmd
	}

	return m, cmd
}

func (m model) View() string {
	var s string

	switch m.currentStep {
	case stepDetectProject:
		s = titleStyle.Render("Detecting YOLO project...") + "\n\n"
		s += docStyle.Render("Please wait...")

	case stepReinitChoice:
		s = titleStyle.Render("YOLO project detected!") + "\n\n"
		s += fmt.Sprintf("Project: %s\n", m.existingConfig.ProjectName)
		s += fmt.Sprintf("AI Provider: %s\n", m.existingConfig.AIProvider)
		s += fmt.Sprintf("Custom Prompts: %v\n\n", m.existingConfig.CustomPrompts)
		s += "What would you like to do?\n\n"
		s += m.configList.View()

	case stepConfirmDelete:
		s = titleStyle.Render("⚠️  Confirm Project Reinitialization") + "\n\n"
		s += "This will delete all YOLO-related files and start fresh.\n"
		s += fmt.Sprintf("Type '%s' to confirm:\n\n", m.existingConfig.ProjectName)
		s += m.textInput.View()

	case stepReconfigureChoice:
		s = titleStyle.Render("Select what to reconfigure:") + "\n\n"
		s += m.configList.View()

	case stepProjectName:
		s = titleStyle.Render("What's your project name?") + "\n\n"
		s += m.textInput.View() + "\n\n"
		s += docStyle.Render("Press enter to confirm")

	case stepProjectPath:
		s = titleStyle.Render("Where should we create the project?") + "\n\n"
		s += m.textInput.View() + "\n\n"
		s += docStyle.Render("Press enter to confirm (default: current directory)")

	case stepGitChoice:
		s = titleStyle.Render("Do you want to use Git?") + "\n\n"
		s += "Git helps track changes and collaborate with others\n\n"
		s += docStyle.Render("(y/N)")

	case stepFolderStructure:
		s = titleStyle.Render("Select folders to include:") + "\n\n"
		s += "(space to toggle, enter to confirm)\n\n"
		s += m.folderList.View()

	case stepAIProvider:
		s = titleStyle.Render("Select your AI provider:") + "\n\n"
		s += m.aiList.View()

	case stepCustomPrompts:
		s = titleStyle.Render("Do you want to customize AI prompts?") + "\n\n"
		s += "This allows you to tailor how the AI assists you\n\n"
		s += docStyle.Render("(y/N)")

	case stepConfirm:
		s = titleStyle.Render("Review your choices:") + "\n\n"
		s += fmt.Sprintf("Project Name: %s\n", m.options.ProjectName)
		s += fmt.Sprintf("Location: %s\n", m.options.ProjectPath)
		s += fmt.Sprintf("Use Git: %v\n", m.options.UseGit)
		s += fmt.Sprintf("AI Provider: %s\n", m.options.AIProvider)
		s += fmt.Sprintf("Custom Prompts: %v\n", m.options.CustomPrompts)
		s += "\nFolders:\n"
		for _, listItem := range m.folderList.Items() {
			i := listItem.(item)
			if i.selected {
				s += fmt.Sprintf("- %s\n", i.title)
			}
		}
		s += "\nPress enter to create project"
	}

	return s
}

func InitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [project-name]",
		Short: "🎉 Start a new project adventure!",
		Long: `🌈 Let's create something amazing together! 

This friendly wizard will help you:
1. 📝 Create a new project or update an existing one
2. 🏠 Set up your project structure
3. 🎨 Configure your settings
4. 🤖 Update templates and documentation
5. 🎯 Keep everything in sync

Options:
  --force    Start fresh (careful: this replaces existing setup)
  --update   Update existing project with new features and templates
  --path     Choose a special spot for your project

Examples:
  yolo init my-awesome-blog
  yolo init "My Cool App"
  yolo init --update  # Update existing project
  yolo init --force   # Start fresh`,
		RunE: runInit,
	}

	cmd.Flags().BoolP("force", "f", false, "✨ Start fresh (careful: this replaces existing setup)")
	cmd.Flags().BoolP("update", "u", false, "🔄 Update existing project with new features")
	cmd.Flags().StringP("path", "p", "", "🏠 Choose a special spot for your project")
	
	return cmd
}

func runInit(cmd *cobra.Command, args []string) error {
	fmt.Println("🎈 Welcome to Your Project Adventure! 🎈")
	
	// Get flags
	force, _ := cmd.Flags().GetBool("force")
	update, _ := cmd.Flags().GetBool("update")
	path, _ := cmd.Flags().GetString("path")
	
	// If path is not specified, use current directory
	if path == "" {
		currentDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		path = currentDir
	}
	
	// Check for existing project
	existingConfig, err := detectExistingProject(path)
	if err != nil {
		return fmt.Errorf("failed to check for existing project: %w", err)
	}

	if existingConfig != nil {
		if force {
			fmt.Println("🔄 Reinitializing project (force mode)...")
			if err := cleanupExistingProject(path); err != nil {
				return fmt.Errorf("failed to cleanup existing project: %w", err)
			}
		} else if update {
			fmt.Println("🔄 Updating existing project...")
			if err := updateExistingProject(path, existingConfig); err != nil {
				return fmt.Errorf("failed to update project: %w", err)
			}
			fmt.Println("\n✨ Project updated successfully!")
			return nil
		} else {
			return fmt.Errorf("directory already contains YOLO files. Use --force to reinitialize or --update to update")
		}
	}
	
	// Change to target directory
	if err := os.Chdir(path); err != nil {
		return fmt.Errorf("failed to change to directory %s: %w", path, err)
	}
	
	// Initialize project
	if err := core.InitializeProject(); err != nil {
		return fmt.Errorf("failed to initialize project: %w", err)
	}
	
	fmt.Println("\n🌟 Success! Your project is ready for amazing things!")
	fmt.Println("\n💡 What's next?")
	fmt.Println("1. Try 'yolo status' to see how things are going")
	fmt.Println("2. Use 'yolo epic' to add your big ideas")
	fmt.Println("3. Run 'yolo graph' to see your project in 3D!")
	
	return nil
}

func updateExistingProject(path string, existingConfig *ProjectConfig) error {
	// Create any missing directories
	newDirs := []string{
		"yolo/relationships", // Added in newer versions
		"yolo/settings",     // Added in newer versions
	}
	
	for _, dir := range newDirs {
		fullPath := filepath.Join(path, dir)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			if err := os.MkdirAll(fullPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", dir, err)
			}
			fmt.Printf("✨ Added new directory: %s\n", dir)
		}
	}

	// Update templates and documentation
	updates := map[string]string{
		"LLM_INSTRUCTIONS.md": templates.LLMInstructionsTemplate, // New AI guidelines
		"yolo/settings/prompts.yml": templates.PromptsTemplate,   // New prompts
	}

	for file, template := range updates {
		fullPath := filepath.Join(path, file)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			content := os.Expand(template, func(key string) string {
				if key == "date" {
					return time.Now().Format("2006-01-02")
				}
				return ""
			})

			if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
				return fmt.Errorf("failed to create file %s: %w", file, err)
			}
			fmt.Printf("✨ Added new file: %s\n", file)
		}
	}

	// Update existing files with new sections
	if err := updateChangelogWithNewFeatures(path); err != nil {
		return fmt.Errorf("failed to update changelog: %w", err)
	}

	if err := updateHistoryWithNewFeatures(path); err != nil {
		return fmt.Errorf("failed to update history: %w", err)
	}

	// Update configuration
	existingConfig.UseConventionalCommits = true // New default
	if err := saveProjectConfig(path, existingConfig); err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}

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
		newContent := string(content) + "\n\n### Work in Progress\n- System-wide keyboard shortcuts (F011)\n  - Web interface\n  - Configuration system\n  - Global shortcuts"
		return os.WriteFile(changelogPath, []byte(newContent), 0644)
	}

	return nil
}

func updateHistoryWithNewFeatures(path string) error {
	historyPath := filepath.Join(path, "HISTORY.yml")
	var history map[string]interface{}

	content, err := os.ReadFile(historyPath)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(content, &history); err != nil {
		return err
	}

	// Add new features if they don't exist
	changes, ok := history["changes"].([]interface{})
	if !ok {
		changes = make([]interface{}, 0)
	}

	// Add new feature entry
	newFeature := map[string]interface{}{
		"type": "feature",
		"id":   "F011",
		"name": "System-wide Keyboard Shortcuts",
		"status": "in_progress",
		"components": []string{"shortcuts", "web", "macos"},
	}

	// Check if feature already exists
	exists := false
	for _, change := range changes {
		if c, ok := change.(map[string]interface{}); ok {
			if id, ok := c["id"].(string); ok && id == "F011" {
				exists = true
				break
			}
		}
	}

	if !exists {
		changes = append(changes, newFeature)
		history["changes"] = changes
		
		newContent, err := yaml.Marshal(history)
		if err != nil {
			return err
		}

		return os.WriteFile(historyPath, newContent, 0644)
	}

	return nil
}

func saveProjectConfig(path string, config *ProjectConfig) error {
	configPath := filepath.Join(path, "yolo", "settings", "config.yml")
	content, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, content, 0644)
}

func cleanupExistingProject(path string) error {
	// Remove YOLO directory
	if err := os.RemoveAll(filepath.Join(path, "yolo")); err != nil {
		return fmt.Errorf("failed to remove YOLO directory: %w", err)
	}

	// Remove base files
	files := []string{
		"CHANGELOG.md",
		"HISTORY.yml",
		"STRATEGY.md",
		"WISHES.md",
		"LLM_INSTRUCTIONS.md",
	}

	for _, file := range files {
		if err := os.Remove(filepath.Join(path, file)); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove %s: %w", file, err)
		}
	}

	return nil
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

	// Create settings
	settings := map[string]interface{}{
		"project_name": opts.ProjectName,
		"ai_provider": opts.AIProvider,
		"use_conventional_commits": opts.UseConventionalCommits,
		"custom_prompts": opts.CustomPrompts,
		"folder_structure": opts.FolderStructure,
	}

	settingsPath := fmt.Sprintf("%s/yolo/settings/config.yml", opts.ProjectPath)
	if err := os.MkdirAll(fmt.Sprintf("%s/yolo/settings", opts.ProjectPath), 0755); err != nil {
		return fmt.Errorf("failed to create settings directory: %w", err)
	}

	// Save settings to YAML file
	settingsData, err := yaml.Marshal(settings)
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	if err := os.WriteFile(settingsPath, settingsData, 0644); err != nil {
		return fmt.Errorf("failed to write settings file: %w", err)
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