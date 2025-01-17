# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- [F011] System-wide keyboard shortcuts feature
  - Web interface for configuring shortcuts
  - JSON-based configuration persistence
  - WebSocket for real-time updates
  - Support for modifier keys (⌘, ⌃, ⌥, ⇧)
  - Command mapping (prompt, graph, commit)

### In Progress
- macOS daemon for global shortcut capture
  - NSGlobalShortcutMonitor integration
  - Accessibility permissions handling
  - System-wide shortcut registration

### Fixed
- Key recording handling for special keys and modifiers
- WebSocket reconnection on disconnect
- Configuration file JSON formatting

## [1.0.0] - 2024-01-16

### Added
- Initial release with core features
- Project visualization
- AI integration
- Git support 