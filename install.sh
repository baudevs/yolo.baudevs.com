#!/bin/bash

# YOLO CLI Installation Script
# This script helps install YOLO CLI and its dependencies

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Required Go version
GO_VERSION="1.21"

# Select personality level
select_personality() {
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
}

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to install Homebrew on macOS
install_homebrew() {
    echo -e "${YELLOW}Installing Homebrew...${NC}"
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
}

# Function to install Go
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

# Function to install Git
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

echo -e "${GREEN}🚀 Welcome to YOLO CLI installer!${NC}"
select_personality

# Check if Git is installed
if ! command_exists git; then
    echo -e "${YELLOW}Git is not installed. Installing...${NC}"
    install_git
else
    echo -e "${GREEN}✅ Git is already installed!${NC}"
fi

# Check if Go is installed
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

# Create temporary directory
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

# Clone and build YOLO
echo -e "${YELLOW}Cloning YOLO CLI repository...${NC}"
git clone https://github.com/baudevs/yolo.baudevs.com.git
cd yolo.baudevs.com

echo -e "${YELLOW}Building YOLO CLI...${NC}"
go mod download
go mod tidy
go build -o yolo cmd/yolo/main.go

echo -e "${YELLOW}Installing YOLO CLI...${NC}"
chmod +x yolo
sudo mv yolo /usr/local/bin/

# Clean up
cd - > /dev/null
rm -rf "$TMP_DIR"

# Final message based on personality
case $(cat "$HOME/.config/yolo/personality") in
    1) echo -e "${GREEN}🎉 Installation complete! Your development environment has been successfully optimized!${NC}" ;;
    2) echo -e "${GREEN}🎉 Done! Try not to break anything important... or do, I'm not your boss!${NC}" ;;
    3) echo -e "${GREEN}🎉 BOOM! We're in business! Time to write some legendary code, you beautiful disaster!${NC}" ;;
esac

echo -e "${YELLOW}To get started, run:${NC}"
echo -e "  yolo init"
echo -e "  yolo ai configure -p openai -k your_api_key"
