#!/bin/bash

# Test script for multiple folder watching functionality
# This script tests the new -paths flag functionality

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Testing multiple folder watching functionality...${NC}"

# Create test directories
TEST_DIR1="/tmp/watch-fs-test1"
TEST_DIR2="/tmp/watch-fs-test2"
TEST_DIR3="/tmp/watch-fs-test3"

# Clean up any existing test directories
rm -rf "$TEST_DIR1" "$TEST_DIR2" "$TEST_DIR3"

# Create test directories
mkdir -p "$TEST_DIR1"
mkdir -p "$TEST_DIR2"
mkdir -p "$TEST_DIR3"

echo -e "${GREEN}Created test directories:${NC}"
echo "  - $TEST_DIR1"
echo "  - $TEST_DIR2"
echo "  - $TEST_DIR3"

# Test 1: Test with single path (backward compatibility)
echo -e "\n${YELLOW}Test 1: Single path (backward compatibility)${NC}"
if ./watch-fs -path "$TEST_DIR1" -version; then
    echo -e "${GREEN}✓ Single path flag works${NC}"
else
    echo -e "${RED}✗ Single path flag failed${NC}"
    exit 1
fi

# Test 2: Test with multiple paths
echo -e "\n${YELLOW}Test 2: Multiple paths${NC}"
PATHS="$TEST_DIR1,$TEST_DIR2,$TEST_DIR3"
if ./watch-fs -paths "$PATHS" -version; then
    echo -e "${GREEN}✓ Multiple paths flag works${NC}"
else
    echo -e "${RED}✗ Multiple paths flag failed${NC}"
    exit 1
fi

# Test 3: Test with spaces in paths
echo -e "\n${YELLOW}Test 3: Paths with spaces${NC}"
PATHS_WITH_SPACES="$TEST_DIR1, $TEST_DIR2 , $TEST_DIR3"
if ./watch-fs -paths "$PATHS_WITH_SPACES" -version; then
    echo -e "${GREEN}✓ Paths with spaces work${NC}"
else
    echo -e "${RED}✗ Paths with spaces failed${NC}"
    exit 1
fi

# Test 4: Test error handling - invalid directory
echo -e "\n${YELLOW}Test 4: Error handling - invalid directory${NC}"
if ./watch-fs -paths "/nonexistent/dir1,/nonexistent/dir2" 2>&1 | grep -q "Invalid directory"; then
    echo -e "${GREEN}✓ Error handling works for invalid directories${NC}"
else
    echo -e "${RED}✗ Error handling failed for invalid directories${NC}"
    exit 1
fi

# Test 5: Test error handling - no paths provided
echo -e "\n${YELLOW}Test 5: Error handling - no paths provided${NC}"
if ./watch-fs 2>&1 | grep -q "either -path or -paths flag is required"; then
    echo -e "${GREEN}✓ Error handling works for missing paths${NC}"
else
    echo -e "${RED}✗ Error handling failed for missing paths${NC}"
    exit 1
fi

# Test 6: Test help message
echo -e "\n${YELLOW}Test 6: Help message${NC}"
if ./watch-fs -h 2>&1 | grep -q "paths" && ./watch-fs -h 2>&1 | grep -q "Comma-separated list"; then
    echo -e "${GREEN}✓ Help message includes paths flag${NC}"
else
    echo -e "${RED}✗ Help message missing paths flag${NC}"
    exit 1
fi

# Clean up
echo -e "\n${YELLOW}Cleaning up test directories...${NC}"
rm -rf "$TEST_DIR1" "$TEST_DIR2" "$TEST_DIR3"

echo -e "\n${GREEN}All tests passed! Multiple folder watching functionality is working correctly.${NC}" 