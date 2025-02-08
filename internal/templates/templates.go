package templates

// LLMInstructionsTemplate provides guidelines for AI/LLM interaction
var LLMInstructionsTemplate = `# AI/LLM Instructions

## Overview

This project follows the YOLO methodology for documentation and project management.
When assisting with this project, please follow these guidelines:

## Core Principles

1. Never delete information
   - Mark as deprecated instead
   - Keep historical context
   - Document reasons for changes

2. Maintain Relationships
   - Link tasks to epics
   - Connect related features
   - Document dependencies
   - Track cross-component impacts

3. Update Documentation
   - Record all changes in CHANGELOG.md
   - Update HISTORY.yml with context
   - Keep README.md current
   - Document decisions in STRATEGY.md

## Status Indicators

Use these emojis consistently:
- ✅ Implemented
- 🏗️ In Progress
- 📅 Planned
- ⚠️ Blocked
- 🔄 Under Review
- 📌 Deprecated

## File Structure

@yolo/
├── @epics/       # Strategic initiatives
├── @features/    # Feature specifications
├── @tasks/       # Implementation tasks
├── @relationships/ # Cross-component links
└── @settings/    # Project configuration

## Version Control

1. Use Conventional Commits
   - feat: New features
   - fix: Bug fixes
   - docs: Documentation
   - style: Formatting
   - refactor: Code restructuring
   - test: Testing
   - chore: Maintenance

2. Follow Semantic Versioning
   - MAJOR: Breaking changes
   - MINOR: New features
   - PATCH: Bug fixes

## Best Practices

1. Always check related files
2. Update all documentation
3. Preserve historical context
4. Link components properly
5. Use proper status indicators
6. Follow naming conventions
7. Add clear descriptions

Remember: The goal is to maintain a clear, traceable history of the project's evolution.`

// PromptsTemplate provides default AI prompts configuration
var PromptsTemplate = `# YOLO AI Prompts

version: 1.0.0
date: ${date}

prompts:
  - name: create-epic
    description: Create a new epic
    template: |
      Create a new epic following YOLO methodology:
      - Title: ${title}
      - Description: ${description}
      - Success Metrics
      - Timeline
      - Related Features
      - Status Indicators

  - name: add-feature
    description: Add a feature to an epic
    template: |
      Add a feature to epic ${epic}:
      - Feature Name: ${name}
      - Description
      - Implementation Tasks
      - Dependencies
      - Status Updates

  - name: create-task
    description: Create implementation task
    template: |
      Create task for feature ${feature}:
      - Task Name: ${name}
      - Implementation Details
      - Acceptance Criteria
      - Dependencies
      - Progress Tracking

  - name: update-status
    description: Update component status
    template: |
      Update status for ${component}:
      - Current Status
      - Progress Details
      - Blockers/Issues
      - Next Steps
      - Related Updates

  - name: document-decision
    description: Record architectural decision
    template: |
      Document decision for ${topic}:
      - Context
      - Considered Options
      - Decision
      - Consequences
      - Related Components

  - name: conventional-commit
    description: Generate commit message
    template: |
      Create conventional commit for ${changes}:
      - Type (feat/fix/etc)
      - Scope
      - Description
      - Breaking Changes
      - Related Issues`

// HistoryTemplate provides the initial HISTORY.yml structure
var HistoryTemplate = `version: 1.0.0
date: ${date}
changes:
  - type: feature
    id: F001
    name: "Project Initialization"
    description: "Initial project setup with YOLO methodology"
    status: implemented
    components:
      - init
      - core
      - templates
    impact: "Established foundation for project documentation"

migrations:
  - version: 1.0.0
    date: ${date}
    changes:
      - "Initial project structure"
      - "Basic documentation setup"
      - "Core file templates"
    rollback:
      - "Remove YOLO directory"
      - "Delete configuration files"`

// ChangelogTemplate provides the initial CHANGELOG.md structure
var ChangelogTemplate = `# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - ${date}

### Added
- Initial project structure
- Basic documentation framework
- Core file templates
- YOLO methodology implementation`

// ReadmeTemplate provides the initial README.md structure
var ReadmeTemplate = `# ${project_name}

Project created with YOLO methodology

## Overview

This project follows the YOLO methodology for documentation and project management.
See YOLO.md in the root directory for full methodology details.

## Quick Start

1. New features go in 'features/'
2. Implementation tasks go in 'tasks/'
3. Strategic initiatives go in 'epics/'
4. Never delete, only mark as deprecated
5. Always link tasks to epics
6. Keep history in CHANGELOG.md and HISTORY.yml`

// YoloReadmeTemplate provides the README.md for the yolo directory
var YoloReadmeTemplate = `# YOLO Documentation

## Project Structure

### Epics
- Strategic initiatives
- Large-scale features
- Project milestones

### Features
- Detailed specifications
- Implementation plans
- Status tracking

### Tasks
- Implementation details
- Progress tracking
- Dependencies

### Relationships
- Component connections
- Dependencies
- Impact analysis

## Documentation Guidelines

1. Never delete information
2. Always add context
3. Link related items
4. Track all changes
5. Use proper status indicators
6. Follow naming conventions`

// StrategyTemplate provides the initial STRATEGY.md structure
var StrategyTemplate = `# Project Strategy

## Vision

Outline the long-term vision and goals for the project.

## Architecture

Document key architectural decisions and their context.

## Roadmap

Plan the project's evolution and major milestones.`

// WishesTemplate provides the initial WISHES.md structure
var WishesTemplate = `# Project Wishes

## Future Improvements

Document desired enhancements and improvements.

## Feature Ideas

Capture potential features and their value proposition.

## Technical Debt

Track areas that need refactoring or improvement.` 