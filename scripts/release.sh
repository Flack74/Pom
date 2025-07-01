#!/bin/bash

# Pom Release Script
# Usage: ./scripts/release.sh [version]

set -e

VERSION=${1:-$(date +"%Y.%m.%d")}
REPO="Flack74/pom"

echo "🚀 Starting release process for version v$VERSION"

# Check if we're in the right directory
if [ ! -f "go.mod" ] || [ ! -d "cmd" ]; then
    echo "❌ Error: Run this script from the project root directory"
    exit 1
fi

# Check if git is clean
if [ -n "$(git status --porcelain)" ]; then
    echo "❌ Error: Git working directory is not clean"
    echo "Please commit or stash your changes first"
    exit 1
fi

# Check if version tag already exists
if git tag -l | grep -q "^v$VERSION$"; then
    echo "❌ Error: Version v$VERSION already exists"
    exit 1
fi

echo "📝 Updating version information..."

# Update version in files if needed
# (Add version updates here if you have version files)

echo "🔨 Building and testing..."

# Build web UI
echo "Building web UI..."
cd web/frontend
npm install
npm run build
cd ../..

# Test build
echo "Testing Go build..."
go build -o pom .
./pom --help > /dev/null

echo "✅ Build successful"

# Run tests if they exist
if [ -f "go.mod" ] && go list ./... | grep -q test; then
    echo "🧪 Running tests..."
    go test ./...
fi

echo "📋 Generating changelog..."

# Get commits since last tag
LAST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
if [ -n "$LAST_TAG" ]; then
    CHANGELOG=$(git log --pretty=format:"- %s" $LAST_TAG..HEAD)
else
    CHANGELOG=$(git log --pretty=format:"- %s" --max-count=10)
fi

echo "📦 Creating release..."

# Create and push tag
git tag -a "v$VERSION" -m "Release v$VERSION"
git push origin "v$VERSION"

echo "⏳ Waiting for GitHub Actions to build artifacts..."
echo "   Check: https://github.com/$REPO/actions"

# Wait a bit for the workflow to start
sleep 10

echo "🎉 Release v$VERSION created successfully!"
echo ""
echo "Next steps:"
echo "1. Wait for GitHub Actions to complete building artifacts"
echo "2. Go to: https://github.com/$REPO/releases/tag/v$VERSION"
echo "3. Edit the release notes if needed"
echo "4. Publish the release"
echo ""
echo "📋 Suggested release notes:"
echo "## What's New in v$VERSION"
echo ""
echo "$CHANGELOG"
echo ""
echo "## Installation"
echo ""
echo "### Arch Linux"
echo "\`\`\`bash"
echo "yay -S pom"
echo "\`\`\`"
echo ""
echo "### Other Platforms"
echo "Download the appropriate binary for your platform from the assets below."
echo ""
echo "## Full Changelog"
echo "**Full Changelog**: https://github.com/$REPO/compare/$LAST_TAG...v$VERSION"

# Clean up
rm -f pom

echo ""
echo "✨ Release process completed!"