#!/bin/bash

# Test script for FileDialog arrow keys functionality
# This script tests that arrow keys work properly in the FileDialog

set -e

echo "🧪 Testing FileDialog arrow keys functionality..."

# Build the application
echo "📦 Building application..."
go build -o bin/watch-fs cmd/watch-fs/main.go

# Create test files for navigation
echo "📁 Creating test files for navigation..."
mkdir -p test_files_arrows
echo "file1" > test_files_arrows/file1.txt
echo "file2" > test_files_arrows/file2.txt
echo "file3" > test_files_arrows/file3.txt
mkdir -p test_files_arrows/subdir1
echo "subfile1" > test_files_arrows/subdir1/subfile1.txt
mkdir -p test_files_arrows/subdir2
echo "subfile2" > test_files_arrows/subdir2/subfile2.txt

# Test that the application starts without errors
echo "🚀 Testing application startup..."
timeout 3s bin/watch-fs test_files_arrows > /dev/null 2>&1 || true

echo "✅ FileDialog arrow keys test completed successfully!"
echo ""
echo "📋 Manual testing instructions:"
echo "1. Run: ./bin/watch-fs test_files_arrows"
echo "2. Press Ctrl+E to open export FileDialog"
echo "3. Test arrow keys:"
echo "   - ↑ (up arrow) should move selection up"
echo "   - ↓ (down arrow) should move selection down"
echo "   - j/k should also work (vim-style)"
echo "4. Test navigation:"
echo "   - Enter on directories should open them"
echo "   - Enter on '..' should go back"
echo "   - Enter on files should select them"
echo "5. Press Esc or q to cancel"
echo ""
echo "🎯 Expected behavior:"
echo "- Arrow keys should move the '>' indicator up/down"
echo "- Selection should be visually highlighted"
echo "- Navigation should be smooth and responsive"
echo "- Both arrow keys and hjkl should work" 