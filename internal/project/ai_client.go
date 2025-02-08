package project

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/baudevs/yolo.baudevs.com/internal/config"
	"github.com/baudevs/yolo.baudevs.com/internal/types"
	"github.com/sashabaranov/go-openai"
)

// AIClient represents an AI-powered client
type AIClient struct {
	client *openai.Client
	model  string
}

// NewAIClient creates a new AI client
func NewAIClient(cfg *config.Config) (*AIClient, error) {
	client := openai.NewClient(cfg.OpenAI.APIKey)
	return &AIClient{
		client: client,
		model:  "gpt-4-turbo-preview",
	}, nil
}

// GenerateProjectName generates a project name based on the description
func (c *AIClient) GenerateProjectName(description string) (string, string, error) {
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "You are a project naming assistant. Generate a short, memorable project name based on the description.",
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: description,
		},
	}

	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    c.model,
			Messages: messages,
		},
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate project name: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", "", fmt.Errorf("no response from AI")
	}

	return resp.Choices[0].Message.Content, description, nil
}

// GenerateProjectPlan generates a project plan based on the description
func (c *AIClient) GenerateProjectPlan(description string) (*types.Project, error) {
	schema := &JSONSchemaDefinition{
		Type: JSONSchemaTypeObject,
		Properties: map[string]JSONSchemaDefinition{
			"name": {
				Type:        JSONSchemaTypeString,
				Description: "The name of the project",
			},
			"description": {
				Type:        JSONSchemaTypeString,
				Description: "The description of the project",
			},
			"epics": {
				Type:        JSONSchemaTypeArray,
				Description: "List of epics in the project",
				Items: &JSONSchemaDefinition{
					Type: JSONSchemaTypeObject,
					Properties: map[string]JSONSchemaDefinition{
						"id": {
							Type:        JSONSchemaTypeString,
							Description: "Unique identifier for the epic (e.g., E001)",
						},
						"name": {
							Type:        JSONSchemaTypeString,
							Description: "Name of the epic",
						},
						"description": {
							Type:        JSONSchemaTypeString,
							Description: "Detailed description of the epic",
						},
						"status": {
							Type:        JSONSchemaTypeString,
							Description: "Current status of the epic",
							Enum:        []string{"planned", "in-progress", "completed"},
						},
						"features": {
							Type:        JSONSchemaTypeArray,
							Description: "List of features in this epic",
							Items: &JSONSchemaDefinition{
								Type: JSONSchemaTypeObject,
								Properties: map[string]JSONSchemaDefinition{
									"id": {
										Type:        JSONSchemaTypeString,
										Description: "Unique identifier for the feature (e.g., F001)",
									},
									"name": {
										Type:        JSONSchemaTypeString,
										Description: "Name of the feature",
									},
									"description": {
										Type:        JSONSchemaTypeString,
										Description: "Detailed description of the feature",
									},
									"status": {
										Type:        JSONSchemaTypeString,
										Description: "Current status of the feature",
										Enum:        []string{"planned", "in-progress", "completed"},
									},
									"tasks": {
										Type:        JSONSchemaTypeArray,
										Description: "List of tasks in this feature",
										Items: &JSONSchemaDefinition{
											Type: JSONSchemaTypeObject,
											Properties: map[string]JSONSchemaDefinition{
												"id": {
													Type:        JSONSchemaTypeString,
													Description: "Unique identifier for the task (e.g., T001)",
												},
												"name": {
													Type:        JSONSchemaTypeString,
													Description: "Name of the task",
												},
												"description": {
													Type:        JSONSchemaTypeString,
													Description: "Detailed description of the task",
												},
												"status": {
													Type:        JSONSchemaTypeString,
													Description: "Current status of the task",
													Enum:        []string{"planned", "in-progress", "completed"},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	schemaJSON, err := json.Marshal(schema)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal schema: %w", err)
	}

	functions := []openai.FunctionDefinition{
		{
			Name:        "generate_project_plan",
			Parameters:  json.RawMessage(schemaJSON),
			Description: "Generate a complete project plan with epics, features, and tasks",
		},
	}

	messages := []openai.ChatCompletionMessage{
		{
			Role: openai.ChatMessageRoleSystem,
			Content: `You are a project planning assistant that creates detailed project plans.
Break down projects into:
1. Epics (major features/components)
2. Features (specific functionalities within epics)
3. Tasks (individual units of work within features)

Generate IDs following these patterns:
- Epics: E001, E002, etc.
- Features: F001, F002, etc.
- Tasks: T001, T002, etc.

All items should start in "planned" status.

IMPORTANT: Always generate a complete project structure with at least one epic, feature, and task.`,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: fmt.Sprintf("Create a project plan for this description:\n\n%s", description),
		},
	}

	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:     c.model,
			Messages:  messages,
			Functions: functions,
			FunctionCall: &openai.FunctionCall{
				Name: "generate_project_plan",
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate project plan: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from AI")
	}

	choice := resp.Choices[0]
	if choice.Message.FunctionCall == nil {
		return nil, fmt.Errorf("no function call in response")
	}

	// Log the response for debugging
	fmt.Println("\nAI Response:")
	fmt.Println("------------------------")
	fmt.Println(choice.Message.FunctionCall.Arguments)
	fmt.Println("------------------------")

	var project types.Project
	if err := json.Unmarshal([]byte(choice.Message.FunctionCall.Arguments), &project); err != nil {
		// If parsing fails, create a default project structure
		project = types.Project{
			Name:        "New Project",
			Description: description,
			Epics: []types.Epic{
				{
					ID:          "E001",
					Name:        "Project Setup",
					Description: "Initial project setup and configuration",
					Status:      "planned",
					Features: []types.Feature{
						{
							ID:          "F001",
							Name:        "Basic Infrastructure",
							Description: "Set up basic project infrastructure and dependencies",
							Status:      "planned",
							Tasks: []types.Task{
								{
									ID:          "T001",
									Name:        "Project Initialization",
									Description: "Initialize project structure and configuration",
									Status:      "planned",
								},
							},
						},
					},
				},
			},
		}
	}

	// Validate minimum requirements
	if project.Name == "" {
		project.Name = "New Project"
	}
	if project.Description == "" {
		project.Description = description
	}
	if len(project.Epics) == 0 {
		// Add default epic structure
		project.Epics = []types.Epic{
			{
				ID:          "E001",
				Name:        "Project Setup",
				Description: "Initial project setup and configuration",
				Status:      "planned",
				Features: []types.Feature{
					{
						ID:          "F001",
						Name:        "Basic Infrastructure",
						Description: "Set up basic project infrastructure and dependencies",
						Status:      "planned",
						Tasks: []types.Task{
							{
								ID:          "T001",
								Name:        "Project Initialization",
								Description: "Initialize project structure and configuration",
								Status:      "planned",
							},
						},
					},
				},
			},
		}
	}

	return &project, nil
}

// EnhanceDescription enhances a project description with more details
func (c *AIClient) EnhanceDescription(description string) (string, error) {
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "You are a project description enhancer. Add more details and clarity to the project description.",
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: fmt.Sprintf("Enhance this project description with more details:\n\n%s", description),
		},
	}

	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    c.model,
			Messages: messages,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to enhance description: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	return resp.Choices[0].Message.Content, nil
}

// GenerateFileContent generates content for a file
func (c *AIClient) GenerateFileContent(fileType string, data interface{}) (string, error) {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal data: %w", err)
	}

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: fmt.Sprintf("You are a file content generator for %s files. Generate detailed content based on the provided data.", fileType),
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: string(dataJSON),
		},
	}

	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    c.model,
			Messages: messages,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate file content: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	return resp.Choices[0].Message.Content, nil
}

// ValidateProjectPlan validates a project plan for completeness and consistency
func (c *AIClient) ValidateProjectPlan(project *types.Project) error {
	projectJSON, err := json.Marshal(project)
	if err != nil {
		return fmt.Errorf("failed to marshal project: %w", err)
	}

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "You are a project plan validator. Check if the project plan is complete and consistent.",
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: fmt.Sprintf("Validate this project plan:\n\n%s", string(projectJSON)),
		},
	}

	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    c.model,
			Messages: messages,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to validate project plan: %w", err)
	}

	if len(resp.Choices) == 0 {
		return fmt.Errorf("no response from AI")
	}

	// Parse validation result
	result := resp.Choices[0].Message.Content
	if strings.Contains(strings.ToLower(result), "invalid") {
		return fmt.Errorf("project plan validation failed: %s", result)
	}

	return nil
}

// ProjectStructure represents the complete project structure
type ProjectStructure struct {
	Project      *types.Project `json:"project"`
	EpicFiles    []FileContent  `json:"epic_files"`
	FeatureFiles []FileContent  `json:"feature_files"`
	TaskFiles    []FileContent  `json:"task_files"`
}

// FileContent represents the content of a file
type FileContent struct {
	ID      string `json:"id"`
	Path    string `json:"path"`
	Content string `json:"content"`
}

// GenerateCompleteProjectStructure generates all project files content
func (c *AIClient) GenerateCompleteProjectStructure(project *types.Project) (*ProjectStructure, error) {
	// Define the schema for the response
	schema := &JSONSchemaDefinition{
		Type: JSONSchemaTypeObject,
		Properties: map[string]JSONSchemaDefinition{
			"project": {
				Type:        JSONSchemaTypeObject,
				Description: "The project details",
				Properties: map[string]JSONSchemaDefinition{
					"name":        {Type: JSONSchemaTypeString},
					"description": {Type: JSONSchemaTypeString},
				},
			},
			"epic_files": {
				Type:        JSONSchemaTypeArray,
				Description: "List of epic files to create",
				Items: &JSONSchemaDefinition{
					Type: JSONSchemaTypeObject,
					Properties: map[string]JSONSchemaDefinition{
						"id":      {Type: JSONSchemaTypeString},
						"path":    {Type: JSONSchemaTypeString},
						"content": {Type: JSONSchemaTypeString},
					},
				},
			},
			"feature_files": {
				Type:        JSONSchemaTypeArray,
				Description: "List of feature files to create",
				Items: &JSONSchemaDefinition{
					Type: JSONSchemaTypeObject,
					Properties: map[string]JSONSchemaDefinition{
						"id":      {Type: JSONSchemaTypeString},
						"path":    {Type: JSONSchemaTypeString},
						"content": {Type: JSONSchemaTypeString},
					},
				},
			},
			"task_files": {
				Type:        JSONSchemaTypeArray,
				Description: "List of task files to create",
				Items: &JSONSchemaDefinition{
					Type: JSONSchemaTypeObject,
					Properties: map[string]JSONSchemaDefinition{
						"id":      {Type: JSONSchemaTypeString},
						"path":    {Type: JSONSchemaTypeString},
						"content": {Type: JSONSchemaTypeString},
					},
				},
			},
		},
		Required: []string{"project", "epic_files", "feature_files", "task_files"},
	}

	schemaJSON, err := json.Marshal(schema)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal schema: %w", err)
	}

	functions := []openai.FunctionDefinition{
		{
			Name:        "generate_project_structure",
			Parameters:  json.RawMessage(schemaJSON),
			Description: "Generate a complete project structure with detailed markdown content for each epic, feature, and task",
		},
	}

	projectJSON, err := json.Marshal(project)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal project: %w", err)
	}

	messages := []openai.ChatCompletionMessage{
		{
			Role: openai.ChatMessageRoleSystem,
			Content: `You are a project structure generator. You will receive a project structure and generate detailed markdown content for each epic, feature, and task.
Each markdown file should follow this structure:

# [Item Title]

## Overview
[A clear and concise description of the item]

## Status
- Current Status: [planned/in-progress/completed]
- Priority: [high/medium/low]
- Timeline: [estimated duration]

## Dependencies
- [List of dependencies with IDs and brief explanations]

## Acceptance Criteria
- [List of specific, measurable criteria that must be met]

## Notes
- [Any additional notes, considerations, or implementation details]

Make sure:
1. Each markdown file follows the exact structure above
2. Content is detailed and specific to the item
3. All sections are properly formatted with correct Markdown syntax
4. Dependencies and related items use correct IDs
5. Status and priority levels are consistent
6. Risk assessments are realistic and include mitigation strategies`,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: fmt.Sprintf("Generate project structure for: %s", string(projectJSON)),
		},
	}

	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:     c.model,
			Messages:  messages,
			Functions: functions,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate project structure: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from AI")
	}

	choice := resp.Choices[0]
	if choice.Message.FunctionCall == nil {
		return nil, fmt.Errorf("no function call in response")
	}

	var structure ProjectStructure
	if err := json.Unmarshal([]byte(choice.Message.FunctionCall.Arguments), &structure); err != nil {
		return nil, fmt.Errorf("failed to parse project structure: %w", err)
	}

	return &structure, nil
}

// GenerateCommitMessage generates a commit message based on the changes
func (c *AIClient) GenerateCommitMessage(changes string) (string, error) {
	messages := []openai.ChatCompletionMessage{
		{
			Role: openai.ChatMessageRoleSystem,
			Content: `You are a commit message generator following the Conventional Commits specification.
Generate commit messages that are:
1. Concise and clear
2. Follow the format: <type>[optional scope]: <description>
3. Types: feat, fix, docs, style, refactor, test, chore
4. Use imperative mood ("add" not "added")
5. No period at the end`,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: fmt.Sprintf("Generate a commit message for these changes:\n\n%s", changes),
		},
	}

	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    c.model,
			Messages: messages,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate commit message: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	return resp.Choices[0].Message.Content, nil
}

// GenerateEpicContent generates markdown content for an epic
func (c *AIClient) GenerateEpicContent(epic types.Epic) (string, error) {
	prompt := fmt.Sprintf(`Create a detailed markdown document for this epic:

Name: %s
Description: %s
Status: %s

The markdown should include:
1. Title with ID and name
2. Description section with comprehensive details
3. Status section with current state
4. Dependencies section if any
5. List of features
6. Any technical considerations
7. Success criteria

Format it as a well-structured markdown document.`, epic.Name, epic.Description, epic.Status)

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "You are a technical documentation writer creating a detailed epic document.",
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: prompt,
		},
	}

	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    c.model,
			Messages: messages,
		},
	)

	if err != nil {
		return "", fmt.Errorf("failed to generate epic content: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	return resp.Choices[0].Message.Content, nil
}

// GenerateFeatureContent generates markdown content for a feature
func (c *AIClient) GenerateFeatureContent(feature types.Feature, parentEpic types.Epic) (string, error) {
	prompt := fmt.Sprintf(`Create a detailed markdown document for this feature:

Name: %s
Description: %s
Status: %s
Parent Epic: [%s] %s

The markdown should include:
1. Title with ID and name
2. Description section with implementation details
3. Status section with current state
4. Parent epic reference
5. List of tasks
6. Technical requirements
7. Success criteria

Format it as a well-structured markdown document.`, feature.Name, feature.Description, feature.Status, parentEpic.ID, parentEpic.Name)

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "You are a technical documentation writer creating a detailed feature document.",
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: prompt,
		},
	}

	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    c.model,
			Messages: messages,
		},
	)

	if err != nil {
		return "", fmt.Errorf("failed to generate feature content: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	return resp.Choices[0].Message.Content, nil
}

// GenerateTaskContent generates markdown content for a task
func (c *AIClient) GenerateTaskContent(task types.Task, parentFeature types.Feature, parentEpic types.Epic) (string, error) {
	prompt := fmt.Sprintf(`Create a detailed markdown document for this task:

Name: %s
Description: %s
Status: %s
Parent Feature: [%s] %s
Parent Epic: [%s] %s

The markdown should include:
1. Title with ID and name
2. Description section with specific implementation details
3. Status section with current state
4. Parent feature and epic references
5. Technical requirements
6. Success criteria
7. Implementation steps

Format it as a well-structured markdown document.`, task.Name, task.Description, task.Status, parentFeature.ID, parentFeature.Name, parentEpic.ID, parentEpic.Name)

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "You are a technical documentation writer creating a detailed task document.",
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: prompt,
		},
	}

	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    c.model,
			Messages: messages,
		},
	)

	if err != nil {
		return "", fmt.Errorf("failed to generate task content: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	return resp.Choices[0].Message.Content, nil
}
