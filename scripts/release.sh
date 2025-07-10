#!/bin/bash

# Release script for watch-fs
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if version is provided
if [ -z "$1" ]; then
    echo -e "${RED}Error: Version number is required${NC}"
    echo "Usage: $0 <version>"
    echo "Example: $0 1.0.0"
    exit 1
fi

VERSION=$1

echo -e "${GREEN}🚀 Starting release process for version $VERSION${NC}"

# Check if we're on main branch
CURRENT_BRANCH=$(git branch --show-current)
if [ "$CURRENT_BRANCH" != "main" ]; then
    echo -e "${RED}Error: Must be on main branch to release${NC}"
    exit 1
fi

# Check if working directory is clean
if [ -n "$(git status --porcelain)" ]; then
    echo -e "${RED}Error: Working directory is not clean${NC}"
    git status
    exit 1
fi

# Run tests
echo -e "${YELLOW}🧪 Running tests...${NC}"
go test -v ./test/...

# Run linter
echo -e "${YELLOW}🔍 Running linter...${NC}"
golangci-lint run

# Build
echo -e "${YELLOW}🔨 Building...${NC}"
make build

# Create and push tag
echo -e "${YELLOW}🏷️  Creating tag v$VERSION...${NC}"
git tag -a "v$VERSION" -m "Release v$VERSION"
git push origin "v$VERSION"

echo -e "${GREEN}✅ Release v$VERSION has been created!${NC}"
echo -e "${GREEN}📦 GitHub Actions will automatically build and create a release${NC}"
echo -e "${GREEN}🔗 Check: https://github.com/pbouamriou/watch-fs/releases${NC}" 