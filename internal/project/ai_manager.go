package project

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/baudevs/yolo.baudevs.com/internal/config"
	"github.com/baudevs/yolo.baudevs.com/internal/types"
)

// AIManager handles AI-powered project management functionality
type AIManager struct {
	client *AIClient
}

// NewAIManager creates a new AI manager instance
func NewAIManager() (*AIManager, error) {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	client, err := NewAIClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create AI client: %w", err)
	}

	return &AIManager{
		client: client,
	}, nil
}

// GetProjectDescription prompts the user for a project description
func (m *AIManager) GetProjectDescription() (string, error) {
	fmt.Println("\nPlease describe your project. You can use formatted text.")
	fmt.Println("Press Cmd+D (Ctrl+D on Linux/Windows) when finished:")

	var description strings.Builder
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		description.WriteString(scanner.Text() + "\n")
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading input: %w", err)
	}

	rawDescription := strings.TrimSpace(description.String())
	if rawDescription == "" {
		return "", fmt.Errorf("project description cannot be empty")
	}

	// Enhance the description using AI
	fmt.Println("\nEnhancing your project description...")
	enhancedDescription, err := m.client.EnhanceDescription(rawDescription)
	if err != nil {
		fmt.Printf("Warning: Failed to enhance description: %v\n", err)
		fmt.Println("Proceeding with original description...")
		return rawDescription, nil
	}

	fmt.Println("\nEnhanced project description:")
	fmt.Println("-----------------------------")
	fmt.Println(enhancedDescription)
	fmt.Println("-----------------------------")

	return enhancedDescription, nil
}

// GetProjectName prompts the user for a project name or generates one
func (m *AIManager) GetProjectName(description string) (string, bool, error) {
	fmt.Print("\nDo you have a name for your project? (yes/no): ")
	var response string
	fmt.Scanln(&response)

	if strings.ToLower(response) == "yes" {
		fmt.Print("Please enter your project name: ")
		var name string
		fmt.Scanln(&name)
		return name, true, nil
	}

	// Generate a project name using AI
	fmt.Println("\nGenerating project name suggestions...")
	generatedName, enhancedDescription, err := m.client.GenerateProjectName(description)
	if err != nil {
		fmt.Printf("Warning: Failed to generate project name: %v\n", err)
		fmt.Println("Please enter your preferred project name: ")
		var name string
		fmt.Scanln(&name)
		return name, true, nil
	}

	// Show the enhanced description
	if enhancedDescription != description {
		fmt.Println("\nEnhanced project description:")
		fmt.Println(enhancedDescription)
	}

	fmt.Printf("\nBased on your description, I suggest: %s\n", generatedName)
	fmt.Print("Would you like to use this name? (yes/no): ")
	fmt.Scanln(&response)

	if strings.ToLower(response) == "yes" {
		return generatedName, false, nil
	}

	fmt.Print("Please enter your preferred project name: ")
	var customName string
	fmt.Scanln(&customName)
	return customName, true, nil
}

// GenerateProjectPlan generates a complete project plan using AI
func (m *AIManager) GenerateProjectPlan(description string) (*types.Project, error) {
	fmt.Println("\nGenerating project plan...")
	project, err := m.client.GenerateProjectPlan(description)
	if err != nil {
		return nil, fmt.Errorf("failed to generate project plan: %w", err)
	}

	fmt.Println("\nProject plan generated successfully!")
	fmt.Printf("✓ %d epics created\n", len(project.Epics))

	var featureCount, taskCount int
	for _, epic := range project.Epics {
		featureCount += len(epic.Features)
		for _, feature := range epic.Features {
			taskCount += len(feature.Tasks)
		}
	}

	fmt.Printf("✓ %d features defined\n", featureCount)
	fmt.Printf("✓ %d tasks identified\n", taskCount)

	return project, nil
}

// SaveProjectStructure saves the project structure to disk
func (m *AIManager) SaveProjectStructure(project *types.Project) error {
	fmt.Println("\n🔄 Generating project structure...")

	// Create project directories
	fmt.Println("\n📁 Creating project directories...")
	dirs := []string{
		"yolo",
		"yolo/epics",
		"yolo/features",
		"yolo/tasks",
	}

	var dirErrors []string
	for _, dir := range dirs {
		fmt.Printf("  Creating %s/...", dir)
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Println(" ❌")
			dirErrors = append(dirErrors, fmt.Sprintf("failed to create directory %s: %v", dir, err))
		} else {
			fmt.Println(" ✓")
		}
	}

	// Save project overview
	fmt.Println("\n📝 Creating project overview...")
	if err := m.saveProjectOverview(project); err != nil {
		fmt.Printf("❌ Failed to create project overview: %v\n", err)
	} else {
		fmt.Println("✓ Created README.md")
	}

	// Track failed files
	var failedFiles []string

	// Create epic files
	fmt.Println("\n📑 Creating epic files...")
	for _, epic := range project.Epics {
		filePath := filepath.Join("yolo/epics", fmt.Sprintf("%s.md", epic.ID))
		fmt.Printf("  Creating %s...", filePath)

		// Generate content for this epic
		content, err := m.client.GenerateEpicContent(epic)
		if err != nil {
			fmt.Println(" ❌")
			failedFiles = append(failedFiles, fmt.Sprintf("%s (content generation error: %v)", filePath, err))
			continue
		}

		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			fmt.Println(" ❌")
			failedFiles = append(failedFiles, fmt.Sprintf("%s (write error: %v)", filePath, err))
			continue
		}
		fmt.Println(" ✓")

		// Create feature files for this epic
		fmt.Printf("\n📑 Creating feature files for epic %s...\n", epic.ID)
		for _, feature := range epic.Features {
			featurePath := filepath.Join("yolo/features", fmt.Sprintf("%s.md", feature.ID))
			fmt.Printf("  Creating %s...", featurePath)

			// Generate content for this feature
			content, err := m.client.GenerateFeatureContent(feature, epic)
			if err != nil {
				fmt.Println(" ❌")
				failedFiles = append(failedFiles, fmt.Sprintf("%s (content generation error: %v)", featurePath, err))
				continue
			}

			if err := os.WriteFile(featurePath, []byte(content), 0644); err != nil {
				fmt.Println(" ❌")
				failedFiles = append(failedFiles, fmt.Sprintf("%s (write error: %v)", featurePath, err))
				continue
			}
			fmt.Println(" ✓")

			// Create task files for this feature
			fmt.Printf("\n📑 Creating task files for feature %s...\n", feature.ID)
			for _, task := range feature.Tasks {
				taskPath := filepath.Join("yolo/tasks", fmt.Sprintf("%s.md", task.ID))
				fmt.Printf("  Creating %s...", taskPath)

				// Generate content for this task
				content, err := m.client.GenerateTaskContent(task, feature, epic)
				if err != nil {
					fmt.Println(" ❌")
					failedFiles = append(failedFiles, fmt.Sprintf("%s (content generation error: %v)", taskPath, err))
					continue
				}

				if err := os.WriteFile(taskPath, []byte(content), 0644); err != nil {
					fmt.Println(" ❌")
					failedFiles = append(failedFiles, fmt.Sprintf("%s (write error: %v)", taskPath, err))
					continue
				}
				fmt.Println(" ✓")
			}
		}
	}

	// Print summary with more details
	fmt.Printf("\n✨ Project structure creation completed!\n")
	fmt.Printf("📊 Summary:\n")
	fmt.Printf("  ✓ %d epics processed\n", len(project.Epics))

	var totalFeatures, totalTasks int
	for _, epic := range project.Epics {
		totalFeatures += len(epic.Features)
		for _, feature := range epic.Features {
			totalTasks += len(feature.Tasks)
		}
	}

	fmt.Printf("  ✓ %d features processed\n", totalFeatures)
	fmt.Printf("  ✓ %d tasks processed\n", totalTasks)

	// Report any failures
	if len(dirErrors) > 0 || len(failedFiles) > 0 {
		fmt.Println("\n⚠️  Some items failed to be created:")
		if len(dirErrors) > 0 {
			fmt.Println("\nDirectory creation errors:")
			for _, err := range dirErrors {
				fmt.Printf("  - %s\n", err)
			}
		}
		if len(failedFiles) > 0 {
			fmt.Println("\nFile creation errors:")
			for _, err := range failedFiles {
				fmt.Printf("  - %s\n", err)
			}
		}
		fmt.Println("\n💡 You may need to create these files manually or run the command again.")
	}

	return nil
}

// verifyCreatedFiles verifies that all files were created correctly
func verifyCreatedFiles(structure *ProjectStructure) error {
	// Verify all epic files
	for _, epicFile := range structure.EpicFiles {
		if _, err := os.Stat(epicFile.Path); err != nil {
			return fmt.Errorf("epic file not found: %s", epicFile.Path)
		}
	}

	// Verify all feature files
	for _, featureFile := range structure.FeatureFiles {
		if _, err := os.Stat(featureFile.Path); err != nil {
			return fmt.Errorf("feature file not found: %s", featureFile.Path)
		}
	}

	// Verify all task files
	for _, taskFile := range structure.TaskFiles {
		if _, err := os.Stat(taskFile.Path); err != nil {
			return fmt.Errorf("task file not found: %s", taskFile.Path)
		}
	}

	// Verify README.md
	if _, err := os.Stat("README.md"); err != nil {
		return fmt.Errorf("README.md not found")
	}

	return nil
}

// saveProjectOverview saves the project overview file
func (m *AIManager) saveProjectOverview(project *types.Project) error {
	content, err := m.client.GenerateFileContent("project", project)
	if err != nil {
		return fmt.Errorf("failed to generate project overview: %w", err)
	}

	if err := os.WriteFile("README.md", []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write project overview: %w", err)
	}

	return nil
}

// GenerateProjectFiles is now handled by SaveProjectStructure
func (m *AIManager) GenerateProjectFiles(project *types.Project) error {
	// This is now handled by SaveProjectStructure
	return nil
}
