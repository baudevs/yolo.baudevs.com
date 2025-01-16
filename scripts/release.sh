#!/bin/bash

# Exit on error
set -e

# Default version bump type
BUMP_TYPE=${1:-patch}

# Validate bump type
if [[ ! "$BUMP_TYPE" =~ ^(major|minor|patch)$ ]]; then
    echo "Error: Version bump type must be one of: major, minor, patch"
    echo "Usage: $0 [major|minor|patch]"
    exit 1
fi

# Get current version from go.mod
CURRENT_VERSION=$(grep "^version =" go.mod | cut -d'"' -f2)
if [ -z "$CURRENT_VERSION" ]; then
    CURRENT_VERSION="0.0.0"
fi

# Split version into components
IFS='.' read -r MAJOR MINOR PATCH <<< "$CURRENT_VERSION"

# Bump version according to type
case $BUMP_TYPE in
    major)
        MAJOR=$((MAJOR + 1))
        MINOR=0
        PATCH=0
        ;;
    minor)
        MINOR=$((MINOR + 1))
        PATCH=0
        ;;
    patch)
        PATCH=$((PATCH + 1))
        ;;
esac

# Construct new version
NEW_VERSION="$MAJOR.$MINOR.$PATCH"
echo "Bumping version from $CURRENT_VERSION to $NEW_VERSION"

# Update version in go.mod
sed -i '' "s/^version = .*/version = \"$NEW_VERSION\"/" go.mod

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

# Stage changes
git add go.mod CHANGELOG.md HISTORY.yml

# Commit changes
git commit -m "chore(release): bump version to $NEW_VERSION"

# Create and push tag
git tag -a "v$NEW_VERSION" -m "Release version $NEW_VERSION"

echo "Release v$NEW_VERSION prepared!"
echo "Next steps:"
echo "1. Review the changes"
echo "2. Push the commit: git push origin main"
echo "3. Push the tag: git push origin v$NEW_VERSION" 