# YOLO: Your Opinionated Lifecycle Orchestrator

YOLO is an AI-powered development methodology and toolset designed to streamline your development workflow by integrating Large Language Models (LLMs) directly into your development process. This guide covers everything you need to know about YOLO and its CLI tools.

## Table of Contents
- [Core Concepts](#core-concepts)
- [Project Structure](#project-structure)
- [Tools & CLI](#tools--cli)
- [Best Practices](#best-practices)

## Core Concepts

### Version Control (yolo commit)
YOLO enhances your Git workflow with AI-powered commit messages. The `yolo commit` command analyzes your changes and generates meaningful, conventional commit messages that:

- Follow semantic versioning conventions
- Include detailed descriptions of changes
- Link to related tasks and features
- Handle large changes gracefully with summarization

Key features:
- Smart diff analysis
- Conventional commit format
- Support for co-authors and issue references
- Handles large changes through chunking or summarization

### History Tracking (yolo history)
YOLO maintains a comprehensive history of your project's evolution through:

- Automatic CHANGELOG generation
- Detailed history tracking in `history.yaml`
- Integration with epics, features, and tasks
- AI-powered summaries of development periods

The history system helps you:
- Track project progress
- Generate reports
- Understand development patterns
- Maintain documentation automatically

### Workflows (yolo epics, tasks, features)
YOLO implements a hierarchical workflow system:

1. **Epics** (`yolo epic`)
   - Large, strategic initiatives
   - Container for features
   - Long-term planning tool
   - AI-assisted scope definition

2. **Features** (`yolo feature`)
   - Concrete functionality implementations
   - Links to parent epics
   - Contains multiple tasks
   - AI-generated implementation plans

3. **Tasks** (`yolo task`)
   - Specific, actionable items
   - Links to parent features
   - Clear success criteria
   - AI-assisted task breakdown

The workflow system ensures:
- Clear hierarchy and relationships
- Automatic linking between items
- AI-powered planning and organization
- Progress tracking at all levels

### Configuration
YOLO can be configured in two ways:

1. **Manual Configuration**
   ```yaml
   # .yolo/config.yaml
   ai:
     provider: openai
     model: gpt-4
   workflow:
     epic_template: templates/epic.md
     feature_template: templates/feature.md
     task_template: templates/task.md
   ```

2. **CLI Configuration**
   ```bash
   yolo config set ai.provider openai
   yolo config set ai.model gpt-4
   yolo config init --interactive
   ```

## Project Structure

### yolo folder
The `.yolo` directory is the heart of your project:

```
.yolo/
├── epics/        # Epic markdown files
├── features/     # Feature markdown files
├── tasks/        # Task markdown files
├── history.yaml  # Detailed history tracking
├── config.yaml   # YOLO configuration
└── templates/    # Custom templates
```

### Key Files

1. **history.yaml**
   ```yaml
   version: 1
   entries:
     - timestamp: "2025-01-23T15:58:01+01:00"
       type: feature
       id: F001
       summary: "Add user authentication"
       related:
         epic: E001
         tasks: [T001, T002]
   ```

2. **CHANGELOG.md**
   - Automatically generated
   - Conventional commit format
   - Grouped by semantic version
   - Links to related items

3. **README.md**
   - Project overview
   - Setup instructions
   - YOLO workflow guide
   - Important links

4. **WISHES.md**
   - Future improvements
   - Feature requests
   - Technical debt items
   - AI-assisted prioritization

5. **STRATEGY.md**
   - Project roadmap
   - Technical decisions
   - Architecture overview
   - AI-generated insights

6. **LLM_INSTRUCTIONS.md**
   - Custom prompts
   - AI interaction guidelines
   - Project-specific context
   - Best practices

## Tools & CLI

### Command Reference

1. **Version Control**
   ```bash
   yolo commit [-a|--all] [-m|--message] [-s|--summarized]
   # -a: Stage all changes
   # -m: Provide custom message
   # -s: Use summarized diff for large changes
   ```

2. **Workflow Management**
   ```bash
   yolo epic create "Epic description"
   yolo epic list [--status=<status>]
   yolo epic update <epic-id>
   
   yolo feature create "Feature description" [-e|--epic <epic-id>]
   yolo feature list [--epic=<epic-id>]
   yolo feature update <feature-id>
   
   yolo task create "Task description" [-f|--feature <feature-id>]
   yolo task list [--feature=<feature-id>]
   yolo task update <task-id>
   ```

3. **History & Reports**
   ```bash
   yolo history show [--from=<date>] [--to=<date>]
   yolo history report [--type=<report-type>]
   yolo changelog generate [--from=<tag>]
   ```

4. **Configuration**
   ```bash
   yolo config init [--interactive]
   yolo config set <key> <value>
   yolo config get <key>
   yolo config list
   ```

5. **AI Interaction**
   ```bash
   yolo ask "Your question"
   yolo explain <file-or-function>
   yolo suggest [--type=<suggestion-type>]
   ```

## Best Practices

### How to Talk to the LLM

1. **Be Specific**
   ```bash
   # Good
   yolo ask "How should I implement JWT authentication in the login endpoint?"
   
   # Bad
   yolo ask "How do I do auth?"
   ```

2. **Provide Context**
   ```bash
   # Good
   yolo explain login.go --context "We're using JWT for authentication"
   
   # Bad
   yolo explain login.go
   ```

3. **Use the Right Command**
   - `yolo ask` for general questions
   - `yolo explain` for code understanding
   - `yolo suggest` for improvements

### Prompt Management

1. **Custom Prompts**
   - Store in `LLM_INSTRUCTIONS.md`
   - Use consistent formatting
   - Include examples
   - Document context requirements

2. **When to Use Each Prompt**
   - Feature creation: High-level planning
   - Task creation: Specific implementation details
   - Code review: Best practices and improvements
   - Bug fixing: Problem analysis

### General Guidelines

1. **Keep Things Small**
   - Break epics into manageable features
   - Split features into focused tasks
   - Make frequent, small commits
   - Use summarized diffs for large changes

2. **Commit Often**
   - After each logical change
   - When tests pass
   - Before switching tasks
   - After refactoring

3. **Maintain Relationships**
   - Link tasks to features
   - Link features to epics
   - Reference issues in commits
   - Keep documentation updated

4. **Use AI Effectively**
   - Let AI handle repetitive tasks
   - Review AI suggestions carefully
   - Provide feedback to improve results
   - Keep prompts up to date

5. **Documentation**
   - Update docs with code changes
   - Use AI for documentation generation
   - Keep README.md current
   - Document AI interactions

Remember: YOLO is designed to enhance your workflow, not replace your judgment. Use the AI tools to augment your development process while maintaining control over the final decisions and implementation details.
