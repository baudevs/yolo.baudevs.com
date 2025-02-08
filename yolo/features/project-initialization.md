# Project Initialization Feature

## Status: Active
## Version: 0.1.0
## Date Created: 2024-01-15
## Related Epic: [CLI Tool](../epics/cli-tool.md)

## Description
The Project Initialization feature provides automated setup of new projects following the YOLO methodology. It creates a standardized directory structure, generates essential documentation files, and optionally integrates existing Git history into the YOLO format.

## Technical Specifications

### Core Components
1. Project Structure (internal/core/initializer.go)
   - Defines standard directory layout
   - Manages file creation
   - Handles template expansion

2. Git Integration (internal/core/git_parser.go)
   - Parses Git history
   - Converts commits to YOLO format
   - Updates history and changelog

### Directory Structure
```
project/
├── CHANGELOG.md
├── HISTORY.yml
├── README.md
├── STRATEGY.md
├── WISHES.md
├── LLM_INSTRUCTIONS.md
└── yolo/
    ├── README.md
    ├── epics/
    ├── features/
    ├── tasks/
    ├── relationships/
    └── settings/
        └── prompts.yml
```

### Templates
- History Template: YAML structure for version tracking
- Changelog Template: Markdown format for changes
- README Template: Project overview
- Strategy Template: Project strategy documentation
- Wishes Template: Project aspirations
- LLM Instructions Template: AI interaction guidelines

### Implementation Details
```go
type ProjectStructure struct {
    Directories []string
    Files       map[string]string
}

func InitializeProject() error {
    // Create directories
    for _, dir := range structure.Directories {
        os.MkdirAll(dir, 0755)
    }

    // Create files from templates
    for path, template := range structure.Files {
        content := os.Expand(template, expandVars)
        os.WriteFile(path, []byte(content), 0644)
    }

    // Optional: Convert Git history
    if versions, err := ParseGitHistory(); err == nil {
        updateHistoryAndChangelog(versions)
    }
}
```

## Impact Analysis
### Current Impact
- Standardizes project setup
- Ensures consistent documentation
- Facilitates YOLO methodology adoption
- Preserves Git history in YOLO format

### Future Impact
- Will support custom templates
- Will enable project-specific configurations
- Will support multiple VCS integrations
- Will add project scaffolding options

## Relationships
### Parent Epic
- [CLI Tool](../epics/cli-tool.md)

### Related Features
- [Active] CLI Framework
- [Active] Prompt Management
- [Planned] Epic Management
- [Planned] Feature Management

### Tasks
- [Implemented] Create directory structure
- [Implemented] Generate documentation files
- [Implemented] Implement Git history parsing
- [Implemented] Add template system
- [Planned] Add custom template support
- [Planned] Implement project configs
- [Planned] Add scaffolding templates

## Notes
- Templates are customizable via settings
- Git history conversion is optional
- All files use standard formats (YAML, Markdown)
- Directory structure follows best practices 