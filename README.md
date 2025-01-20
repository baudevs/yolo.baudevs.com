# 🚀 YOLO CLI - Your Project Companion

YOLO (You Only Live Once) is your cheeky, over-caffeinated project management sidekick. It’s fun, it’s sassy, and it’s built to make managing your AI coding projects ridiculously easy. Whether you’re a seasoned developer or an LLM just trying to vibe, YOLO’s got your back. 🤖

## ✨ Features

- 🎯 **Project Organization**: Keep your chaos in check with a clean, standardized structure.
- 🤖 **AI Love**: Works seamlessly with LLMs. Yes, even the moody ones.
- 📊 **Visual Progress**: 3D visualizations that’ll make you go “Oooooh!”
- 🚀 **Blazing Fast**: Because who has time to wait? Certainly not YOLO.
- 🔄 **Git Smarts**: Effortlessly manage commits like a version control wizard.

## 🛠️ Installation

### Quick Install (Recommended)

Run our installation script that handles everything for you:

```bash
curl -fsSL https://raw.githubusercontent.com/baudevs/yolo.baudevs.com/main/install.sh | bash
```

The installer will:
- Install Git if needed
- Install Go if needed
- Clone the repository
- Build and install YOLO CLI
- Guide you through the setup

### Manual Installation

If you prefer to install manually:

1. **Install Prerequisites**
   ```bash
   # macOS
   brew install git go

   # Ubuntu/Debian
   sudo apt-get update && sudo apt-get install git golang-go

   # Check installations
   git --version
   go version  # Should be 1.21 or later
   ```

2. **Build YOLO**
   ```bash
   # Clone repository
   git clone https://github.com/baudevs/yolo.baudevs.com.git
   cd yolo.baudevs.com

   # Build and install
   make
   make install
   ```

3. **Configure YOLO**
   ```bash
   # Initialize a project
   yolo init

   # Configure AI provider
   yolo ai configure -p openai -k your_api_key
   ```

## 📚 Project Structure

YOLO sets up your project like a boss. Here’s what it’ll look like:

```
your-project/
├── CHANGELOG.md      # Project changes and versions
├── HISTORY.yml       # Complete historical record
├── README.md         # Project overview
├── STRATEGY.md       # Project strategy and goals
├── WISHES.md         # Future improvements
├── LLM_INSTRUCTIONS.md  # AI/LLM guidelines
└── yolo/
    ├── epics/       # Strategic initiatives
    ├── features/    # Feature specifications
    ├── tasks/       # Implementation tasks
    ├── relationships/ # Cross-component links
    └── settings/    # Project configuration
```

---

## 🎮 YOLO Commands

YOLO CLI isn’t just a pretty face. Here’s what it can do:

- `yolo init`: Set up a new project adventure.
- `yolo prompt`: Get AI-powered assistance (because thinking is overrated).
- `yolo graph`: Visualize project relationships like a boss.
- `yolo commit`: Create AI-powered conventional commits with style.
- `yolo shortcuts`: Configure global shortcuts (WIP).
- `yolo ai`: Manage AI providers and settings:
  - `yolo ai configure`: Set up your AI providers
  - `yolo ai list`: View configured providers
  - `yolo ai test`: Test your AI setup

### 🤖 AI Configuration

YOLO supports multiple AI providers to power its features:

```bash
# Configure OpenAI
yolo ai configure -p openai -k your_api_key

# List providers
yolo ai list

# Test configuration
yolo ai test
```

Supported providers:
- OpenAI (default)
- Anthropic Claude
- Mistral AI

### 🎯 Smart Commit

YOLO's commit command is now powered by AI:

```bash
# Create an AI-powered commit
yolo commit

# Skip remote sync
yolo commit --no-sync

# Force commit with warnings
yolo commit --force
```

Features:
- AI-generated conventional commits
- Automatic staging
- Remote repository syncing
- Smart error handling with AI assistance

---

## ⚙️ Development Setup

For those brave enough to tinker under the hood:

```bash
# Clone the YOLO repository
git clone https://github.com/baudevs/yolo-cli
cd yolo-cli

# Install dependencies
go mod download

# Build the project
go build -o bin/yolo cmd/yolo/main.go

# Run tests
go test ./...
```

---

## 🤝 Contributing

Got ideas? Found bugs? Wanna show off your coding chops? Check out our [Contributing Guide](CONTRIBUTING.md). We welcome PRs, memes, and bribes (just kidding, kind of). 🤪

---

## 📜 License

MIT © [BauDevs](https://baudevs.com)

---

Built with 💥, too much caffeine and lots of chaos (you wanted love?), by the BauDevs team and Monoverse. Special shoutouts to:

- [BauDevs](https://baudevs.social)
- [storres3rd](https://github.com/storres3rd)
- [Monoverse](https://monoverse.com)
YOLO responsibly! 🎉
