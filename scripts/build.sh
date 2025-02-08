#!/bin/bash

# Set version from argument or use default
VERSION=${1:-"0.1.0"}
BINARY_NAME="yolo"
PACKAGE_NAME="github.com/baudevs/yolo-cli"

# Create release directory
mkdir -p release

# Build for each platform
PLATFORMS=("windows/amd64" "windows/arm64" "darwin/amd64" "darwin/arm64" "linux/amd64" "linux/arm64")

for platform in "${PLATFORMS[@]}"
do
    # Split platform into OS and architecture
    IFS='/' read -r -a array <<< "$platform"
    GOOS="${array[0]}"
    GOARCH="${array[1]}"
    
    # Set output binary name based on OS
    if [ "$GOOS" = "windows" ]; then
        output_name="$BINARY_NAME.exe"
    else
        output_name="$BINARY_NAME"
    fi

    # Build binary
    echo "Building for $GOOS/$GOARCH..."
    export GOOS=$GOOS
    export GOARCH=$GOARCH
    go build -o "release/${BINARY_NAME}-${GOOS}-${GOARCH}/${output_name}" ./cmd/yolo

    # Create archive
    cd release
    if [ "$GOOS" = "windows" ]; then
        zip -r "${BINARY_NAME}-${VERSION}-${GOOS}-${GOARCH}.zip" "${BINARY_NAME}-${GOOS}-${GOARCH}"
    else
        tar -czf "${BINARY_NAME}-${VERSION}-${GOOS}-${GOARCH}.tar.gz" "${BINARY_NAME}-${GOOS}-${GOARCH}"
    fi
    cd ..

    echo "✓ Done building for $GOOS/$GOARCH"
done

echo "All builds completed!" 