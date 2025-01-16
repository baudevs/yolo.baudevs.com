# YOLO CLI Tool

## Description

YOLO CLI is a powerful command-line tool designed to implement and manage the YOLO (You Observe, Log, and Oversee) methodology. This methodology is specifically designed for effective collaboration between developers and Large Language Models (LLMs) while maintaining a complete and transparent project history.

## Features

### Active Features
- **Project Initialization**: Automatically set up new projects with YOLO methodology structure
- **Documentation Management**: Maintain comprehensive project documentation with historical context
- **Prompt Management**: Access and customize LLM interaction prompts
- **Git Integration**: Convert existing Git history to YOLO format

### Planned Features
- **Epic Management**: Track large-scale project initiatives
- **Feature Management**: Document and track feature development
- **Task Management**: Create and monitor task progress
- **Status Tracking**: Monitor overall project health

## Installation

### Prerequisites
- Go 1.21 or higher
- Git

### Building from Source
```bash
# Clone the repository
git clone https://github.com/baudevs/yolo-cli.git
cd yolo-cli

# Build the binary
go build -o bin/yolo cmd/yolo/main.go

# Optional: Add to PATH
cp bin/yolo /usr/local/bin/
```

### macOS and Linux
```bash
curl -fsSL https://raw.githubusercontent.com/baudevs/yolo-cli/main/scripts/install.sh | bash
```

To install a specific version:
```bash
curl -fsSL https://raw.githubusercontent.com/baudevs/yolo-cli/main/scripts/install.sh | bash -s VERSION
```

### Windows
Using PowerShell:
```powershell
irm https://raw.githubusercontent.com/baudevs/yolo-cli/main/scripts/install.ps1 | iex
```

To install a specific version:
```powershell
$Version="VERSION"; irm https://raw.githubusercontent.com/baudevs/yolo-cli/main/scripts/install.ps1 | iex
```

Alternatively, you can download the appropriate zip file for your architecture (amd64 or arm64) from the [releases page](https://github.com/baudevs/yolo-cli/releases) and extract it to a directory in your PATH.

## Usage

### Initialize a New Project
```bash
yolo init
```

### Work with Prompts
```bash
# Get standard documentation prompt
yolo prompt standard

# Get changelog update prompt
yolo prompt changelog

# Get readme update prompt
yolo prompt readme

# Get epic documentation prompt
yolo prompt epic

# Get feature documentation prompt
yolo prompt feature

# Get task documentation prompt
yolo prompt task

# Get history update prompt
yolo prompt history
```

## Project Structure
```
project/
├── CHANGELOG.md      # Project changelog (SemVer)
├── HISTORY.yml      # Complete change history
├── README.md        # Project documentation
├── STRATEGY.md      # Project strategy
├── WISHES.md        # Project aspirations
├── LLM_INSTRUCTIONS.md  # AI interaction guidelines
└── yolo/
    ├── README.md    # YOLO documentation
    ├── epics/       # Epic documentation
    ├── features/    # Feature documentation
    ├── tasks/       # Task documentation
    ├── relationships/  # Relationship tracking
    └── settings/    # Project settings
        └── prompts.yml  # LLM prompts
```

## YOLO Methodology

The YOLO methodology is built on these key principles:

1. **Complete History**: Never delete information, only mark it as deprecated, superseded, or not relevant
2. **Relationship Tracking**: Maintain clear connections between epics, features, and tasks
3. **LLM Collaboration**: Standardized prompts and documentation for effective AI interaction
4. **Version Control**: Follow SemVer while preserving historical context
5. **Documentation First**: Comprehensive documentation as a core development practice

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (following [Conventional Commits](https://www.conventionalcommits.org/))
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Inspired by the need for better developer-AI collaboration
- Built with [Cobra](https://github.com/spf13/cobra)
- Uses [YAML v3](https://github.com/go-yaml/yaml)