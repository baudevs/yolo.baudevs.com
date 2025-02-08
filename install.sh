#!/bin/bash

# 🚀 YOLO CLI Installation Script
# This script installs or updates YOLO CLI and its dependencies.
#
# Features:
# - Smart detection of existing installations
# - Preserves user configuration during updates
# - Supports development mode for local builds
# - Personality-based feedback and messages
#
# Usage:
# - Fresh install: curl -fsSL https://raw.githubusercontent.com/baudevs/yolo.baudevs.com/main/install.sh | bash
# - Update: Run the script again, it will detect existing installation
# - Development: Run from repo directory for local build

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Required Go version
GO_VERSION="1.21"

# Detect if YOLO is already initialized by checking for config directory
# and essential files like personality settings
is_yolo_initialized() {
    # Check for YOLO config directory and essential files
    [ -d "$HOME/.config/yolo" ] && [ -f "$HOME/.config/yolo/personality" ]
}

# Get the current version of YOLO CLI if installed
# Returns "not_installed" if YOLO is not found
get_current_version() {
    CURRENT_VERSION=$(grep "Version = " internal/version/version.go | cut -d'"' -f2)
    if [ -n "$CURRENT_VERSION" ]; then
        echo "$CURRENT_VERSION"
    else
        echo "unknown"
    fi
}
    

# Compare version strings
# Returns 1 if version1 is greater, -1 if version2 is greater, 0 if equal
compare_versions() {
    if [[ $1 == $2 ]]; then
        echo 0
        return
    fi
    
    local IFS=.
    local i ver1=($1) ver2=($2)
    
    # Fill empty positions in ver1 with zeros
    for ((i=${#ver1[@]}; i<${#ver2[@]}; i++)); do
        ver1[i]=0
    done
    
    # Fill empty positions in ver2 with zeros
    for ((i=${#ver2[@]}; i<${#ver1[@]}; i++)); do
        ver2[i]=0
    done
    
    # Compare version numbers
    for ((i=0; i<${#ver1[@]}; i++)); do
        if ((10#${ver1[i]} > 10#${ver2[i]})); then
            echo 1
            return
        fi
        if ((10#${ver1[i]} < 10#${ver2[i]})); then
            echo -1
            return
        fi
    done
    
    echo 0
}

# Detect if we're running in development mode
# This is true if we're in the repo directory with go.mod
is_dev_mode() {
    # Check if we're in the repo directory or its subdirectories
    git rev-parse --git-dir > /dev/null 2>&1 && [[ -f "go.mod" ]]
}

# Configure or load YOLO's personality
# This function either prompts for personality selection (new install)
# or loads existing personality (update)
select_personality() {
    # Only prompt for personality if not already set
    if [ ! -f "$HOME/.config/yolo/personality" ]; then
        echo -e "${GREEN}Select YOLO's personality level:${NC}"
        echo "1) Clean & Nerdy (Safe for work, still fun)"
        echo "2) Mildly Eccentric (Slightly edgy, occasional sass)"
        echo "3) Unhinged & Funny (Full chaos mode, not for the faint of heart)"
        
        read -p "Enter your choice (1-3) [default: 1]: " personality
        
        # Create YOLO config directory
        CONFIG_DIR="$HOME/.config/yolo"
        mkdir -p "$CONFIG_DIR"
        
        # Save personality to config file
        case $personality in
            2) echo "2" > "$CONFIG_DIR/personality" ;;
            3) echo "3" > "$CONFIG_DIR/personality" ;;
            *) echo "1" > "$CONFIG_DIR/personality" ;;
        esac
        
        # Create shell configuration
        SHELL_TYPE=$(basename "$SHELL")
        SHELL_CONFIG=""
        
        case "$SHELL_TYPE" in
            "zsh")  SHELL_CONFIG="$HOME/.zshenv" ;;
            "bash") SHELL_CONFIG="$HOME/.bash_profile" ;;
            *)      SHELL_CONFIG="$HOME/.profile" ;;
        esac
        
        # Remove any existing YOLO_PERSONALITY export
        if [ -f "$SHELL_CONFIG" ]; then
            sed -i.bak '/export YOLO_PERSONALITY=/d' "$SHELL_CONFIG"
        fi
        
        # Add new personality setting
        echo "export YOLO_PERSONALITY=$(cat "$CONFIG_DIR/personality")" >> "$SHELL_CONFIG"
        
        # Export for current session
        export YOLO_PERSONALITY=$(cat "$CONFIG_DIR/personality")
        
        case $personality in
            1) echo -e "${GREEN}Excellent choice! Let's proceed with scientific precision! 🧪${NC}" ;;
            2) echo -e "${GREEN}Oh, feeling sassy today, are we? Let's do this! 😏${NC}" ;;
            3) echo -e "${GREEN}YOLO MODE ACTIVATED! Hold onto your bits! 🚀${NC}" ;;
        esac
    else
        # Load existing personality
        export YOLO_PERSONALITY=$(cat "$HOME/.config/yolo/personality")
    fi
}

# Create default prompts configuration
create_default_prompts() {
    CONFIG_DIR="$HOME/.config/yolo"
    PROMPTS_FILE="$CONFIG_DIR/prompts.yml"
    
    # Only create if it doesn't exist
    if [ ! -f "$PROMPTS_FILE" ]; then
        cat > "$PROMPTS_FILE" << 'EOL'
messages:
  welcome:
    nerdy_clean: "🚀 Welcome to YOLO CLI - Your Optimal Life Organizer!"
    mildly_rude: "🤘 Sup nerd! Welcome to YOLO - Let's break some stuff!"
    unhinged_funny: "🔥 YOLO CLI in da house! Time to code like nobody's watching!"
  install_start:
    nerdy_clean: "Initiating installation sequence with quantum precision..."
    mildly_rude: "Alright, let's get this party started! Installing stuff..."
    unhinged_funny: "Hold onto your bits! We're about to go full send on this install!"
  install_go:
    nerdy_clean: "Installing Go - The language of gophers 🐹"
    mildly_rude: "Yo, we need Go! Don't ask questions, just let it happen..."
    unhinged_funny: "Time to inject some Go juice into your machine! YEET! 🚀"
  install_git:
    nerdy_clean: "Installing Git - Version control for the win!"
    mildly_rude: "Need Git because apparently you live in 2025... Installing!"
    unhinged_funny: "Installing Git cuz we ain't savages! Time to get version controlled! 🎮"
  install_done:
    nerdy_clean: "Installation complete! Your development environment has been optimized."
    mildly_rude: "Done! Try not to break anything important... or do, I'm not your boss!"
    unhinged_funny: "BOOM! We're in business! Time to write some legendary code, you beautiful disaster!"
  init_start:
    nerdy_clean: "Initializing YOLO project structure with mathematical precision..."
    mildly_rude: "Let's get this party started! Time to YOLO-ify your project..."
    unhinged_funny: "YOLO MODE ENGAGED! Prepare for project transformation! 🚀"
  init_done:
    nerdy_clean: "Project initialized successfully! Ready for optimal productivity."
    mildly_rude: "Project's all set up! Don't mess it up... too much."
    unhinged_funny: "BOOM! Your project just got YOLO'd! Let the chaos begin! 🎉"
  commit_start:
    nerdy_clean: "Analyzing changes with AI precision..."
    mildly_rude: "Let's see what mess you've made this time..."
    unhinged_funny: "Time to let the AI judge your code! No pressure! 😈"
  commit_done:
    nerdy_clean: "Changes committed successfully! Your code is now immortalized."
    mildly_rude: "Alright, your changes are in! Hope you tested them... maybe."
    unhinged_funny: "YEET! Your code is now part of history! No takebacks! 🚀"
EOL
    fi
}

# Check if a command exists in the system
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Install Homebrew package manager on macOS
install_homebrew() {
    echo -e "${YELLOW}Installing Homebrew...${NC}"
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
}

# Install Go programming language
# Supports both macOS (via Homebrew) and Linux (via apt/yum)
install_go() {
    echo -e "${YELLOW}Installing Go...${NC}"
    case "$(uname -s)" in
        Darwin)
            if ! command_exists brew; then
                install_homebrew
            fi
            brew install go
            ;;
        Linux)
            if command_exists apt-get; then
                sudo apt-get update
                sudo apt-get install -y golang-go
            elif command_exists yum; then
                sudo yum install -y golang
            else
                echo -e "${RED}Unsupported Linux distribution. Please install Go manually: https://golang.org/doc/install${NC}"
                exit 1
            fi
            ;;
        *)
            echo -e "${RED}Unsupported operating system. Please install Go manually: https://golang.org/doc/install${NC}"
            exit 1
            ;;
    esac
}

# Install Git version control system
# Supports both macOS (via Homebrew) and Linux (via apt/yum)
install_git() {
    echo -e "${YELLOW}Installing Git...${NC}"
    case "$(uname -s)" in
        Darwin)
            if ! command_exists brew; then
                install_homebrew
            fi
            brew install git
            ;;
        Linux)
            if command_exists apt-get; then
                sudo apt-get update
                sudo apt-get install -y git
            elif command_exists yum; then
                sudo yum install -y git
            else
                echo -e "${RED}Unsupported Linux distribution. Please install Git manually: https://git-scm.com/downloads${NC}"
                exit 1
            fi
            ;;
        *)
            echo -e "${RED}Unsupported operating system. Please install Git manually: https://git-scm.com/downloads${NC}"
            exit 1
            ;;
    esac
}

# Create user's local bin directory if it doesn't exist
ensure_local_bin() {
    LOCAL_BIN="$HOME/.local/bin"
    mkdir -p "$LOCAL_BIN"
    
    # Add to PATH if not already there
    case ":$PATH:" in
        *":$LOCAL_BIN:"*) ;;
        *)
            # Determine shell config file
            SHELL_TYPE=$(basename "$SHELL")
            case "$SHELL_TYPE" in
                "zsh")  SHELL_CONFIG="$HOME/.zshrc" ;;
                "bash") 
                    if [[ "$OSTYPE" == "darwin"* ]]; then
                        SHELL_CONFIG="$HOME/.bash_profile"
                    else
                        SHELL_CONFIG="$HOME/.bashrc"
                    fi
                    ;;
                *)      SHELL_CONFIG="$HOME/.profile" ;;
            esac
            
            # Add to PATH in shell config if not already there
            if ! grep -q "export PATH=\"\$HOME/.local/bin:\$PATH\"" "$SHELL_CONFIG" 2>/dev/null; then
                echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$SHELL_CONFIG"
                export PATH="$HOME/.local/bin:$PATH"
            fi
            ;;
    esac
}

# Main installation flow
# Check if YOLO is already installed and handle accordingly
LATEST_VERSION="1.1.0" # This should be updated with each release
CURRENT_VERSION=$(get_current_version)

if [ "$CURRENT_VERSION" != "not_installed" ] && [ "$CURRENT_VERSION" != "unknown" ]; then
    VERSION_COMPARE=$(compare_versions "$LATEST_VERSION" "$CURRENT_VERSION")
    if [ "$VERSION_COMPARE" -le 0 ]; then
        echo -e "${GREEN}✨ YOLO CLI is already up to date (version $CURRENT_VERSION)${NC}"
        exit 0
    else
        echo -e "${YELLOW}🚀 Updating YOLO CLI from version $CURRENT_VERSION to $LATEST_VERSION...${NC}"
    fi
else
    echo -e "${YELLOW}🎉 Installing YOLO CLI version $LATEST_VERSION...${NC}"
fi

# Select or load personality
select_personality

# Create default prompts if needed
create_default_prompts

# Check and install dependencies
if ! command_exists git; then
    echo -e "${YELLOW}Git is not installed. Installing...${NC}"
    install_git
else
    echo -e "${GREEN}✅ Git is already installed!${NC}"
fi

if ! command_exists go; then
    echo -e "${YELLOW}Go is not installed. Installing...${NC}"
    install_go
else
    CURRENT_GO_VERSION=$(go version | cut -d' ' -f3 | sed 's/go//')
    if [ "$CURRENT_GO_VERSION" \< "$GO_VERSION" ]; then
        echo -e "${YELLOW}Go version $CURRENT_GO_VERSION is older than required $GO_VERSION. Updating...${NC}"
        install_go
    else
        echo -e "${GREEN}✅ Go $GO_VERSION or later is already installed!${NC}"
    fi
fi

# Build and install YOLO
echo -e "${YELLOW}Building YOLO CLI...${NC}"

if is_dev_mode; then
    echo -e "${GREEN}📦 Development mode detected - using local files${NC}"
    # We're already in the repo directory
    BUILD_DIR="$(pwd)"
else
    echo -e "${YELLOW}Cloning YOLO CLI repository...${NC}"
    # Create temporary directory and clone
    BUILD_DIR=$(mktemp -d)
    cd "$BUILD_DIR"
    git clone https://github.com/baudevs/yolo.baudevs.com.git
    cd yolo.baudevs.com
fi

# Build YOLO
go mod download
go mod tidy
go build -o yolo-cli cmd/yolo/main.go

echo -e "${YELLOW}Installing YOLO CLI...${NC}"
# Ensure local bin directory exists and is in PATH
ensure_local_bin

# Install YOLO to user's local bin
chmod +x yolo-cli
if [ -f "$HOME/.local/bin/yolo" ]; then
    rm -f "$HOME/.local/bin/yolo"
elif [ -d "$HOME/.local/bin/yolo" ]; then
    rm -rf "$HOME/.local/bin/yolo"
fi
cp yolo-cli "$HOME/.local/bin/yolo"
rm yolo-cli

# Clean up if not in dev mode
if [ "$BUILD_DIR" != "$(pwd)" ]; then
    cd - > /dev/null
    rm -rf "$BUILD_DIR"
fi

# Final message based on personality and installation type
case $(cat "$HOME/.config/yolo/personality") in
    1) echo -e "${GREEN}🎉 YOLO CLI has been $(if [ "$CURRENT_VERSION" != "not_installed" ]; then echo "updated"; else echo "installed"; fi) successfully! Your development environment has been optimized!${NC}" ;;
    2) echo -e "${GREEN}🎉 Done! $(if [ "$CURRENT_VERSION" != "not_installed" ]; then echo "Updated all your stuff"; else echo "Try not to break anything important... or do, I'm not your boss"; fi)!${NC}" ;;
    3) echo -e "${GREEN}🎉 BOOM! $(if [ "$CURRENT_VERSION" != "not_installed" ]; then echo "Fresh code injected"; else echo "Time to write some legendary code"; fi), you beautiful disaster!${NC}" ;;
esac

echo -e "${YELLOW}To get started, run:${NC}"
if [ "$CURRENT_VERSION" == "not_installed" ]; then
    echo -e "  yolo init"
fi
echo -e "  yolo ai configure -k your_api_key \n  yolo init \n"
