# Changelog

All notable changes to this project will be documented in this file.

## [0.1.3] - 2024-01-31

### Changed
- Migrated license management from direct Stripe integration to backend API
- Updated license activation to use new verification endpoint
- Improved credit tracking and validation through backend API

### Added
- New `license checkout` command for purchasing licenses
- Support for different license packages (basic, pro, enterprise)
- Enhanced license status display with better formatting

### Removed
- Direct Stripe integration code
- Local credit syncing (now handled by backend)
