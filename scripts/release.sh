#!/bin/bash

# YOLO Release Process Documentation
# ================================
#
# This script handles the release process for YOLO, supporting stable releases and pre-releases (alpha/beta).
#
# Version Format
# -------------
# - Stable: X.Y.Z (e.g., 1.0.0)
# - Beta: X.Y.Z-beta.N (e.g., 1.0.0-beta.1)
# - Alpha: X.Y.Z-alpha.N (e.g., 1.0.0-alpha.1)
#
# Usage
# -----
# 1. Normal version bumping:
#    ./scripts/release.sh [major|minor|patch] [stable|beta|alpha]
#    
#    Examples:
#    - ./scripts/release.sh patch         # 1.0.0 -> 1.0.1
#    - ./scripts/release.sh minor beta    # 1.0.0 -> 1.1.0-beta.1
#    - ./scripts/release.sh major alpha   # 1.0.0 -> 2.0.0-alpha.1
#
# 2. Force specific version:
#    ./scripts/release.sh --version X.Y.Z[-beta.N|-alpha.N]
#    
#    Examples:
#    - ./scripts/release.sh --version 1.2.3        # Force version 1.2.3
#    - ./scripts/release.sh --version 1.2.3-beta.2 # Force version 1.2.3-beta.2
#
# Release Process
# -------------
# 1. Prerequisites:
#    - Must be on main branch
#    - Working directory must be clean
#    - Latest changes must be pulled from main
#
# 2. Automated Steps (handled by this script):
#    - Creates release branch (release/vX.Y.Z)
#    - Updates version in:
#      * internal/version/version.go
#      * install.sh
#    - Updates CHANGELOG.md
#    - Updates HISTORY.yml
#    - Runs build and tests
#    - Creates git commit
#    - Creates version tag
#
# 3. Manual Steps (after script runs):
#    a) Review the changes:
#       git diff main
#
#    b) Push release branch:
#       git push origin release/vX.Y.Z
#
#    c) Create Pull Request:
#       - Title: chore(release): bump version to X.Y.Z
#       - Base: main
#       - Compare: release/vX.Y.Z
#
#    d) After PR approval and merge:
#       git checkout main
#       git pull origin main
#       git push origin vX.Y.Z
#
#    e) Create GitHub Release:
#       - For stable: Create normal release
#       - For alpha/beta: Create pre-release
#       - Use CHANGELOG.md content for description
#
# Reverting a Release
# -----------------
# If something goes wrong, you can revert with:
#   git tag -d vX.Y.Z
#   git reset --hard HEAD^
#   git checkout main
#
# Files Modified
# -------------
# 1. internal/version/version.go - Main version number
# 2. install.sh - Installation script version
# 3. CHANGELOG.md - Release notes and changes
# 4. HISTORY.yml - Release history
#
# Version Bumping Rules
# -------------------
# 1. Major (X.0.0): Breaking changes
#    - Resets minor and patch to 0
#    - For pre-release: adds -alpha.1 or -beta.1
#
# 2. Minor (X.Y.0): New features, backward compatible
#    - Resets patch to 0
#    - For pre-release: adds -alpha.1 or -beta.1
#
# 3. Patch (X.Y.Z): Bug fixes, backward compatible
#    - Increments patch number
#    - For pre-release: adds -alpha.1 or -beta.1
#
# 4. Pre-release (-alpha.N or -beta.N):
#    - If same type exists: increments N
#    - If different type: starts at 1

# Exit on error
set -e

# Parse arguments
FORCE_VERSION=""
while [[ $# -gt 0 ]]; do
    case $1 in
        --version)
            FORCE_VERSION="$2"
            shift 2
            ;;
        *)
            if [ -z "$BUMP_TYPE" ]; then
                BUMP_TYPE="$1"
            elif [ -z "$RELEASE_TYPE" ]; then
                RELEASE_TYPE="$1"
            else
                echo "Error: Unknown argument $1"
                echo "Usage: $0 [major|minor|patch] [stable|beta|alpha] [--version X.X.X]"
                exit 1
            fi
            shift
            ;;
    esac
done

# Set defaults if not provided
BUMP_TYPE=${BUMP_TYPE:-patch}
RELEASE_TYPE=${RELEASE_TYPE:-stable}

# Validate bump type if not using force version
if [ -z "$FORCE_VERSION" ] && [[ ! "$BUMP_TYPE" =~ ^(major|minor|patch)$ ]]; then
    echo "Error: Version bump type must be one of: major, minor, patch"
    echo "Usage: $0 [major|minor|patch] [stable|beta|alpha] [--version X.X.X]"
    exit 1
fi

# Validate release type if not using force version
if [ -z "$FORCE_VERSION" ] && [[ ! "$RELEASE_TYPE" =~ ^(stable|beta|alpha)$ ]]; then
    echo "Error: Release type must be one of: stable, beta, alpha"
    echo "Usage: $0 [major|minor|patch] [stable|beta|alpha] [--version X.X.X]"
    exit 1
fi

# Validate force version format if provided
if [ -n "$FORCE_VERSION" ]; then
    if ! [[ "$FORCE_VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+(-alpha\.[0-9]+|-beta\.[0-9]+)?$ ]]; then
        echo "Error: Force version must be in format X.X.X or X.X.X-[alpha|beta].X"
        echo "Examples: 1.0.0, 1.0.0-beta.1, 1.0.0-alpha.2"
        exit 1
    fi
    
    # Determine release type from forced version
    if [[ "$FORCE_VERSION" =~ -beta\.[0-9]+$ ]]; then
        RELEASE_TYPE="beta"
    elif [[ "$FORCE_VERSION" =~ -alpha\.[0-9]+$ ]]; then
        RELEASE_TYPE="alpha"
    else
        RELEASE_TYPE="stable"
    fi
fi

# Ensure we're on main branch and up to date
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [ "$CURRENT_BRANCH" != "main" ]; then
    echo "Error: Must be on main branch to create a release"
    exit 1
fi

# Ensure working directory is clean
if [ -n "$(git status --porcelain)" ]; then
    echo "Error: Working directory is not clean. Please commit or stash changes."
    exit 1
fi

# Pull latest changes
echo "Pulling latest changes from main..."
git pull origin main

# Get current version from version.go
CURRENT_VERSION=$(grep "Version = " internal/version/version.go | cut -d'"' -f2)
if [ -z "$CURRENT_VERSION" ]; then
    CURRENT_VERSION="0.0.0"
fi

if [ -n "$FORCE_VERSION" ]; then
    NEW_VERSION="$FORCE_VERSION"
    echo "Using forced version: $NEW_VERSION (${RELEASE_TYPE} release)"
else
    # Extract base version and pre-release info
    BASE_VERSION=$(echo "$CURRENT_VERSION" | sed -E 's/(-alpha|-beta)\.[0-9]+$//')
    PRE_RELEASE=$(echo "$CURRENT_VERSION" | grep -Eo '(-alpha|-beta)\.[0-9]+$' || echo "")
    PRE_RELEASE_NUM=0

    if [ -n "$PRE_RELEASE" ]; then
        PRE_RELEASE_NUM=$(echo "$PRE_RELEASE" | grep -Eo '[0-9]+$')
    fi

    # Split version into components
    IFS='.' read -r MAJOR MINOR PATCH <<< "$BASE_VERSION"

    # Bump version according to type
    case $BUMP_TYPE in
        major)
            if [ "$RELEASE_TYPE" = "stable" ]; then
                MAJOR=$((MAJOR + 1))
                MINOR=0
                PATCH=0
                PRE_RELEASE=""
            else
                MAJOR=$((MAJOR + 1))
                MINOR=0
                PATCH=0
                PRE_RELEASE_NUM=1
            fi
            ;;
        minor)
            if [ "$RELEASE_TYPE" = "stable" ]; then
                MINOR=$((MINOR + 1))
                PATCH=0
                PRE_RELEASE=""
            else
                MINOR=$((MINOR + 1))
                PATCH=0
                PRE_RELEASE_NUM=1
            fi
            ;;
        patch)
            if [ "$RELEASE_TYPE" = "stable" ]; then
                PATCH=$((PATCH + 1))
                PRE_RELEASE=""
            else
                PATCH=$((PATCH + 1))
                PRE_RELEASE_NUM=1
            fi
            ;;
    esac

    # Construct new version
    NEW_VERSION="$MAJOR.$MINOR.$PATCH"

    # Add pre-release suffix if not stable
    case $RELEASE_TYPE in
        beta)
            # If already a beta, increment the number
            if [[ "$CURRENT_VERSION" =~ -beta\.[0-9]+$ ]]; then
                PRE_RELEASE_NUM=$((PRE_RELEASE_NUM + 1))
            fi
            NEW_VERSION="$NEW_VERSION-beta.$PRE_RELEASE_NUM"
            ;;
        alpha)
            # If already an alpha, increment the number
            if [[ "$CURRENT_VERSION" =~ -alpha\.[0-9]+$ ]]; then
                PRE_RELEASE_NUM=$((PRE_RELEASE_NUM + 1))
            fi
            NEW_VERSION="$NEW_VERSION-alpha.$PRE_RELEASE_NUM"
            ;;
    esac

    echo "Bumping version from $CURRENT_VERSION to $NEW_VERSION"
fi

# Create release branch
BRANCH_NAME="release/v$NEW_VERSION"
echo "Creating release branch $BRANCH_NAME..."
git checkout -b "$BRANCH_NAME"

# Update version in version.go
sed -i '' "s/Version = \".*\"/Version = \"$NEW_VERSION\"/" internal/version/version.go

# Update version in install.sh
sed -i '' "s/LATEST_VERSION=\".*\"/LATEST_VERSION=\"$NEW_VERSION\"/" install.sh

# Update CHANGELOG.md
DATE=$(date +%Y-%m-%d)
sed -i '' "s/\[Unreleased\]/[$NEW_VERSION] - $DATE/" CHANGELOG.md

# Update HISTORY.yml
cat >> HISTORY.yml << EOL

version: $NEW_VERSION
date: $DATE
history:
  - type: chore
    scope: release
    subject: "bump version to $NEW_VERSION"
    body: "Release version $NEW_VERSION"
EOL

# Build and test
echo "Building and testing..."
make build test

# Stage changes
git add internal/version/version.go install.sh CHANGELOG.md HISTORY.yml

# Commit changes
git commit -m "chore(release): bump version to $NEW_VERSION"

# Create tag
git tag -a "v$NEW_VERSION" -m "Release version $NEW_VERSION"

echo "✨ Release v$NEW_VERSION prepared successfully!"
echo
echo "Release Checklist:"
echo "1. Review the changes:"
echo "   git diff main"
echo
echo "2. Push the release branch:"
echo "   git push origin $BRANCH_NAME"
echo
echo "3. Create a Pull Request:"
echo "   - Title: 'chore(release): bump version to $NEW_VERSION'"
echo "   - Base: main"
echo "   - Compare: $BRANCH_NAME"
echo
echo "4. After PR is approved and merged:"
echo "   git checkout main"
echo "   git pull origin main"
echo "   git push origin v$NEW_VERSION"
echo
if [ "$RELEASE_TYPE" = "stable" ]; then
    echo "5. Create a GitHub Release:"
    echo "   - Tag: v$NEW_VERSION"
    echo "   - Title: Release v$NEW_VERSION"
    echo "   - Description: Copy the changes from CHANGELOG.md"
    echo
    echo "6. The new version will be available after CI/CD completes"
else
    echo "5. Create a GitHub Pre-release:"
    echo "   - Tag: v$NEW_VERSION"
    echo "   - Title: Pre-release v$NEW_VERSION"
    echo "   - Description: Copy the changes from CHANGELOG.md"
    echo "   - Check 'This is a pre-release' option"
    echo
    echo "6. The new pre-release version will be available after CI/CD completes"
fi
echo
echo "To revert this release if needed:"
echo "   git tag -d v$NEW_VERSION"
echo "   git reset --hard HEAD^"
echo "   git checkout main"