# 🚀 YOLO CLI - Your Project Companion

YOLO (You Only Live Once) is a revolutionary project management methodology and CLI tool that makes development fun, organized, and AI-friendly! 

## ✨ Features

- 🎯 **Project Organization**: Create and manage projects with a clear, standardized structure
- 🤖 **AI Integration**: Built-in support for AI-powered assistance and code generation
- 📊 **Visual Progress**: See your project come alive with 3D visualization
- ⌨️ **Global Shortcuts**: Quick access to commands from anywhere (🚧 Work in Progress)
- 📝 **Smart Documentation**: Never lose context with our unique history-preserving approach
- 🔄 **Git Integration**: Seamless version control with conventional commits support

## 🚀 Quick Start

```bash
# Install YOLO CLI
brew tap baudevs/yolo
brew install yolo-cli

# Create a new project
mkdir my-awesome-project
cd my-awesome-project
yolo init

# Or initialize with options
yolo init --path /custom/path
yolo init --force  # Reinitialize existing project
```

## 📚 Project Structure

```
your-project/
├── CHANGELOG.md      # Project changes and versions
├── HISTORY.yml      # Complete historical record
├── README.md        # Project overview
├── STRATEGY.md      # Project strategy and goals
├── WISHES.md        # Future improvements
├── LLM_INSTRUCTIONS.md  # AI/LLM guidelines
└── yolo/
    ├── epics/       # Strategic initiatives
    ├── features/    # Feature specifications
    ├── tasks/       # Implementation tasks
    ├── relationships/ # Cross-component links
    └── settings/    # Project configuration
```

## 🎮 Commands

- `yolo init`: Start a new project adventure
- `yolo prompt`: Get AI-powered assistance
- `yolo graph`: Visualize project relationships
- `yolo commit`: Create conventional commits
- `yolo shortcuts`: Configure global shortcuts (🚧 WIP)

## 🎯 Status

- **Version**: 1.0.0-beta
- **Stability**: Beta
- **Node Support**: v18+
- **Go Version**: 1.21+

### 🚧 Work in Progress

- System-wide keyboard shortcuts
  - Web interface: ✅ Implemented
  - Configuration: ✅ Implemented
  - macOS Daemon: 🏗️ In Progress
  - Linux Support: 📅 Planned

## 🤝 Contributing

We love contributions! Check out our [Contributing Guide](CONTRIBUTING.md) to get started.

### Development Setup

```bash
# Clone the repository
git clone https://github.com/baudevs/yolo-cli
cd yolo-cli

# Install dependencies
go mod download

# Build the project
go build -o bin/yolo cmd/yolo/main.go

# Run tests
go test ./...
```

## 📜 License

MIT © [BauDevs](https://baudevs.com)

---

Built with ❤️ by the BauDevs team
