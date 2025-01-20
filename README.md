# 🚀 YOLO CLI - Your Project Companion

YOLO (You Only Live Once) is your cheeky project management sidekick. Choose your preferred personality level and let YOLO guide you through your coding adventures! 

## ✨ Features

- 🎯 **Project Organization**: Keep your chaos in check with a clean, standardized structure.
- 🤖 **AI Love**: Works seamlessly with LLMs. Yes, even the moody ones.
- 📊 **Visual Progress**: 3D visualizations that'll make you go "Oooooh!"
- 🚀 **Blazing Fast**: Because who has time to wait? Certainly not YOLO.
- 🔄 **Git Smarts**: Effortlessly manage commits like a version control wizard.

## 🎭 Personality Levels

YOLO comes with three distinct personality levels:

1. **Clean & Nerdy** (Default)
   - Safe for work
   - Nerdy humor
   - Professional but fun

2. **Mildly Eccentric**
   - Slightly edgy
   - Occasional sass
   - Work-appropriate but spicy

3. **Unhinged & Funny**
   - Full chaos mode
   - Not for the faint of heart
   - Maximum entertainment

## 🛠️ Installation

### Quick Install (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/baudevs/yolo.baudevs.com/main/install.sh | bash
```

The installer will:
- Detect if YOLO is already installed
- Preserve your existing configuration and preferences
- Update or install components as needed
- Guide you through any additional setup

### Update Flow

When updating an existing YOLO installation:
1. Your personality settings are preserved
2. Configuration files remain untouched
3. Only the CLI binary is updated
4. New features are seamlessly integrated

### Fresh Installation

For new installations, the installer will:
1. Let you choose YOLO's personality
2. Install Git if needed
3. Install Go if needed
4. Set up your configuration
5. Guide you through initial setup

### Manual Installation

1. **Install Prerequisites**
   ```bash
   # macOS
   brew install git go
   
   # Ubuntu/Debian
   sudo apt-get update
   sudo apt-get install git golang-go
   ```

2. **Clone & Build**
   ```bash
   git clone https://github.com/baudevs/yolo.baudevs.com.git
   cd yolo.baudevs.com
   go build -o yolo cmd/yolo/main.go
   ```

3. **Install**
   ```bash
   sudo mv yolo /usr/local/bin/
   ```

## 🎨 Configuration

### Personality Setup

Your personality choice is stored in `~/.config/yolo/personality`:
```
1 = Clean & Nerdy
2 = Mildly Eccentric
3 = Unhinged & Funny
```

To change personality:
```bash
echo "2" > ~/.config/yolo/personality  # For Mildly Eccentric
source ~/.zshenv  # Or your shell's config file
```

### AI Configuration

Configure your AI provider:
```bash
yolo ai configure -p openai -k your_api_key
```

## 🚀 Getting Started

### First Time Setup

```bash
# Initialize YOLO in your project
yolo init

# Configure AI features
yolo ai configure -p openai -k your_api_key

# Start the visualization server
yolo graph
```

### Daily Usage

```bash
# Create a new feature
yolo new feature "Add awesome stuff"

# Let AI analyze your changes
yolo commit

# View your progress
yolo graph
```

## 🎮 Commands

### Core Commands
- `init`: Initialize YOLO in your project
- `new`: Create new features or tasks
- `commit`: Create smart commits with AI
- `graph`: Launch the visualization server

### AI Commands
- `ai configure`: Set up AI providers
- `ai test`: Test your AI configuration

### Utility Commands
- `version`: Show YOLO version
- `update`: Update YOLO CLI
- `help`: Show help message

## 🎯 Development

### Requirements
- Go 1.21 or later
- Git
- An OpenAI API key (for AI features)

### Local Development

1. Clone the repo:
   ```bash
   git clone https://github.com/baudevs/yolo.baudevs.com.git
   cd yolo.baudevs.com
   ```

2. Install dependencies:
   ```bash
   go mod download
   go mod tidy
   ```

3. Build:
   ```bash
   go build -o yolo cmd/yolo/main.go
   ```

4. Run tests:
   ```bash
   go test ./...
   ```

## 🤝 Contributing

1. Fork it
2. Create your feature branch (`git checkout -b feature/amazing`)
3. Commit your changes (`yolo commit`)
4. Push to the branch (`git push origin feature/amazing`)
5. Create a new Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

Built with 💥, too much caffeine and lots of chaos (you wanted love?), by the BauDevs team and Monoverse. Special shoutouts to:

- [BauDevs](https://baudevs.social)
- [storres3rd](https://github.com/storres3rd)
- [Monoverse](https://monoverse.com)
