#!/bin/bash

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Map architecture names
case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

# Set version and binary name
VERSION=${1:-"0.1.0"}
BINARY_NAME="yolo"
INSTALL_DIR="/usr/local/bin"

# Create temporary directory
TMP_DIR=$(mktemp -d)
cd $TMP_DIR

# Download the appropriate release
RELEASE_URL="https://github.com/baudevs/yolo-cli/releases/download/v${VERSION}/${BINARY_NAME}-${VERSION}-${OS}-${ARCH}.tar.gz"
echo "Downloading YOLO CLI v${VERSION} for ${OS}/${ARCH}..."
curl -L -o release.tar.gz $RELEASE_URL

# Extract the archive
tar xzf release.tar.gz

# Install the binary
echo "Installing YOLO CLI to ${INSTALL_DIR}..."
sudo mv "${BINARY_NAME}-${OS}-${ARCH}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
sudo chmod +x "${INSTALL_DIR}/${BINARY_NAME}"

# Clean up
cd - > /dev/null
rm -rf $TMP_DIR

echo "✓ YOLO CLI has been installed successfully!"
echo "Run 'yolo --help' to get started." 