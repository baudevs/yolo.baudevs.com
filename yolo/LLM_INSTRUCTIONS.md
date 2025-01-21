# AI/LLM Instructions

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

Remember: The goal is to maintain a clear, traceable history of the project's evolution.