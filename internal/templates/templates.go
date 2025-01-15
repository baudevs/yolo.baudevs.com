package templates

const (
	HistoryTemplate = `version: 0.1.0
date: ${date}
changes:
  - type: feature
    description: "Initial project setup with YOLO methodology"
    impact: "Established foundation for project documentation"
    files:
      - "README.md"
      - "HISTORY.yml"
      - "CHANGELOG.md"
      - "WISHES.md"
      - "STRATEGY.md"
      - "LLM_INSTRUCTIONS.md"`

	ChangelogTemplate = `# Changelog

All notable changes to this project will be documented in this file.`

	ReadmeTemplate = `# Project Name

## Description

This project follows the YOLO methodology for development and documentation.`

	WishesTemplate = `# Project Wishes

Document your project wishes and aspirations here.`

	StrategyTemplate = `# Project Strategy

Document your project strategy here.`

	YoloReadmeTemplate = `# YOLO Documentation

This directory contains YOLO methodology documentation.`

	PromptsTemplate = `# YOLO Methodology Prompts

# Documentation prompts
standard_documentation: |
  awesome! now document everything according to our yolo methodology in @yolo folder adding the epics, tasks and the relationships between all of them, update the @README @WISHES @STRATEGY and @HISTORY and @CHANGELOG  as we do in the YOLO methodology. Remember: we don't delete anything or lines, we just mark as deprecated, added, implemented, not relevant, anything that changes, so we can keep a complete history of the evolution of our project. @CHANGELOG.md and @README.md maintain standard development practices, like SemVer. Remember we never delete anything, we always add lines, mark as deprecated, implemented, not relevant, updated, improved, created, fixed, deleted, and keep adding context to the project.

update_changelog: |
  Please update the @CHANGELOG.md following the YOLO methodology and SemVer practices. Remember to never delete entries, only add new ones, and mark old ones as deprecated, updated, or superseded when necessary.

update_readme: |
  Please update the @README.md to reflect the latest changes while following YOLO methodology. Remember to preserve historical context by marking outdated sections as deprecated rather than removing them, and add new sections for current information.

epic_documentation: |
  Please document this epic in the @yolo/epics directory following YOLO methodology. Include relationships with features and tasks, impact analysis, and maintain historical context of any changes.

feature_documentation: |
  Please document this feature in the @yolo/features directory following YOLO methodology. Include relationships with epics and tasks, technical specifications, and maintain historical context of any changes.

task_documentation: |
  Please document this task in the @yolo/tasks directory following YOLO methodology. Include relationships with epics and features, implementation details, and maintain historical context of any changes.

update_history: |
  Please update the @HISTORY.yml to reflect these changes following YOLO methodology. Remember to include the type of change, description, impact, and affected files while maintaining the complete history.`

	LLMInstructionsTemplate = `# LLM Instructions

Instructions for LLM interactions with this project.

For standard documentation prompts, check yolo/settings/prompts.yml
Use 'yolo rp' to copy the standard documentation prompt to your clipboard.`
) 