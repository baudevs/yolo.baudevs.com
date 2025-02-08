package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func ExplainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "explain",
		Short: "🎓 Learn about YOLO methodology",
		Long: `🌟 Understand the YOLO (You Only Live Once) methodology in detail.
		
This command provides a comprehensive explanation of:
- Project structure and hierarchy
- File relationships and dependencies
- Best practices and conventions
- AI/LLM integration guidelines`,
		RunE: runExplain,
	}
	return cmd
}

func runExplain(cmd *cobra.Command, args []string) error {
	explanation := `# 🎯 YOLO Methodology Explained

## Core Principles

YOLO is built on three fundamental principles:
1. Never delete, only mark as changed/deprecated
2. Maintain complete historical context
3. Make relationships explicit and traceable

## 📁 Directory Structure (@folders)

@yolo/
├── @epics/       # Strategic initiatives and large-scale features
├── @features/    # Detailed feature specifications
├── @tasks/       # Implementation tasks and subtasks
├── @relationships/ # Cross-component connections
└── @settings/    # Project configuration and preferences

## 📄 Core Files (@files)

- @CHANGELOG.md: Version-based changes (SemVer format)
- @HISTORY.yml: Complete historical record with relationships
- @README.md: Project overview and quick start
- @STRATEGY.md: Long-term goals and architectural decisions
- @WISHES.md: Future improvements and feature requests
- @LLM_INSTRUCTIONS.md: AI/LLM specific guidelines

## 🔄 Relationships

1. Epics (@epics/) contain:
   - Strategic goals
   - Feature groups
   - Success metrics
   - Timeline estimates

2. Features (@features/) link to:
   - Parent epic(s)
   - Implementation tasks
   - Related features
   - Status updates

3. Tasks (@tasks/) connect to:
   - Parent feature
   - Related tasks
   - Implementation details
   - Progress tracking

4. Relationships (@relationships/) track:
   - Epic-to-Epic dependencies
   - Feature-to-Feature links
   - Task dependencies
   - Cross-component impacts

## 📝 Documentation Rules

1. All changes must be recorded in:
   - @CHANGELOG.md for version tracking
   - @HISTORY.yml for complete context
   - Related component files

2. Status updates use emojis:
   - ✅ Implemented
   - 🏗️ In Progress
   - 📅 Planned
   - ⚠️ Blocked
   - 🔄 Under Review
   - 📌 Deprecated

3. Version control follows:
   - Semantic Versioning
   - Conventional Commits
   - Clear change descriptions

## 🤖 LLM Instructions

As an AI assistant, you should:
1. Always preserve historical context
2. Update all related files when making changes
3. Maintain explicit relationships between components
4. Use proper status indicators
5. Follow the folder structure
6. Link tasks to parent features/epics
7. Record all decisions and their context

## 🔍 File Relationships

1. @HISTORY.yml tracks:
   - All changes chronologically
   - Component relationships
   - Migration paths
   - Rollback procedures

2. @CHANGELOG.md contains:
   - Version-based grouping
   - User-facing changes
   - Breaking changes
   - Migration guides

3. Component files include:
   - Links to related components
   - Status indicators
   - Change history
   - Context preservation

## 📊 Version Management

1. Version Format: MAJOR.MINOR.PATCH
   - MAJOR: Breaking changes
   - MINOR: New features
   - PATCH: Bug fixes

2. Each version requires:
   - @CHANGELOG.md update
   - @HISTORY.yml entry
   - Component documentation
   - Relationship updates

## 🎯 Best Practices

1. Never delete information
2. Always link related components
3. Maintain clear status indicators
4. Preserve decision context
5. Update all related files
6. Follow naming conventions
7. Use proper emoji indicators

Remember: After any change, always:
1. Update @HISTORY.yml
2. Update @CHANGELOG.md
3. Update related components
4. Verify relationships
5. Add proper status indicators
6. Preserve historical context`

	fmt.Println(explanation)
	return nil
}
