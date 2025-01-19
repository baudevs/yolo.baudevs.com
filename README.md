# 🚀 YOLO CLI - Your Project Companion

YOLO (You Only Live Once) is your cheeky, over-caffeinated project management sidekick. It’s fun, it’s sassy, and it’s built to make managing your AI coding projects ridiculously easy. Whether you’re a seasoned developer or an LLM just trying to vibe, YOLO’s got your back. 🤖

## ✨ Features

- 🎯 **Project Organization**: Keep your chaos in check with a clean, standardized structure.
- 🤖 **AI Love**: Works seamlessly with LLMs. Yes, even the moody ones.
- 📊 **Visual Progress**: 3D visualizations that’ll make you go “Oooooh!”
- 🚀 **Blazing Fast**: Because who has time to wait? Certainly not YOLO.
- 🔄 **Git Smarts**: Effortlessly manage commits like a version control wizard.

## 🛠️ Installation Instructions

Alright, listen up, champ! Here’s how to install YOLO CLI and unleash your inner coding beast:

### 1. **Download YOLO for macOS** (Sorry, Windows and Linux folks, you’re next on the list. Soon™.)

- Go to the [YOLO Releases Page](https://github.com/baudevs/yolo.baudevs.com/releases/tag/v1.0.0-beta).
- Download the file named `yolo-darwin` (macOS = Darwin, don’t ask, it’s science).

### 2. **Move YOLO to the VIP PATH**

Open the **Terminal** (that’s the hacker-y black box thing) and run these commands:

```bash
# Make YOLO executable
chmod +x ~/Downloads/yolo-darwin

# Move it to your PATH (so YOLO can flex everywhere)
mv ~/Downloads/yolo-darwin /usr/local/bin/yolo
```

### 3. **Party Time!**

Run this command from your project root:

```bash
yolo init
```

🎉 Boom! You’re now officially YOLO-ing like a pro.

---

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
- `yolo commit`: Create conventional commits with style.
- `yolo shortcuts`: Configure global shortcuts (WIP).

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

Built with 💥, bad jokes, and too much caffeine by the BauDevs team and Monoverse. Special shoutouts to:

- [juanda](https://github.com/baudevs.social)
- [storres3rd](https://github.com/storres3rd/playlistsource)

YOLO responsibly! 🎉
