#!/bin/bash

# Function to increment version
increment_version() {
    local version=$1
    local is_beta=$(echo $version | grep -c "beta")
    
    if [ $is_beta -eq 1 ]; then
        # For beta versions (e.g., 1.0.0-beta.3 -> 1.0.0-beta.4)
        local beta_num=$(echo $version | sed 's/.*beta\.//')
        local base_version=$(echo $version | sed 's/-beta\..*//')
        local new_beta_num=$((beta_num + 1))
        echo "${base_version}-beta.${new_beta_num}"
    else
        # For regular versions (e.g., 1.0.0 -> 1.0.1)
        local parts=($(echo $version | tr '.' ' '))
        local major=${parts[0]}
        local minor=${parts[1]}
        local patch=${parts[2]}
        echo "$major.$minor.$((patch + 1))"
    fi
}

# Get current version from version.go
CURRENT_VERSION=$(grep 'Version = ' internal/version/version.go | sed 's/.*Version = "\(.*\)".*/\1/')
echo "Current version: $CURRENT_VERSION"

# Calculate new version
NEW_VERSION=$(increment_version $CURRENT_VERSION)
echo "New version: $NEW_VERSION"

# Confirm with user
read -p "Do you want to proceed with version $NEW_VERSION? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Operation cancelled"
    exit 1
fi

# Create new branch
BRANCH_NAME="release/v$NEW_VERSION"
git checkout -b $BRANCH_NAME
if [ $? -ne 0 ]; then
    echo "Failed to create branch $BRANCH_NAME"
    exit 1
fi
echo "Created branch $BRANCH_NAME"

# Update version.go
sed -i '' "s/Version = \"$CURRENT_VERSION\"/Version = \"$NEW_VERSION\"/" internal/version/version.go
if [ $? -ne 0 ]; then
    echo "Failed to update version.go"
    exit 1
fi
echo "Updated version.go"

# Update install.sh
sed -i '' "s/VERSION=\${1:-\".*\"}/VERSION=\${1:-\"$NEW_VERSION\"}/" scripts/install.sh
if [ $? -ne 0 ]; then
    echo "Failed to update install.sh"
    exit 1
fi
echo "Updated install.sh"

# Show changes
echo -e "\nChanges made:"
git diff

echo -e "\n✨ Version bump completed!"
echo "Next steps:"
echo "1. Review the changes"
echo "2. Commit with: git commit -am 'chore(release): bump version to $NEW_VERSION'"
echo "3. Push with: git push origin $BRANCH_NAME"
echo "4. Create a pull request" 