#!/bin/bash

# Debug script for FileDialog arrow keys issue
# This script helps identify why arrow keys don't work in FileDialog

set -e

echo "🔍 Debugging FileDialog arrow keys issue..."

# Build the application
echo "📦 Building application..."
go build -o bin/watch-fs cmd/watch-fs/main.go

# Create test files
echo "📁 Creating test files..."
mkdir -p test_debug
echo "test1" > test_debug/file1.txt
echo "test2" > test_debug/file2.txt
echo "test3" > test_debug/file3.txt
mkdir -p test_debug/subdir
echo "subtest" > test_debug/subdir/subfile.txt

echo "✅ Debug setup completed!"
echo ""
echo "🔧 Debugging steps:"
echo "1. Run: ./bin/watch-fs test_debug"
echo "2. Press Ctrl+I to open import FileDialog"
echo "3. Try these keys and observe behavior:"
echo "   - ↑ (up arrow) - should move selection up"
echo "   - ↓ (down arrow) - should move selection down"
echo "   - j/k - should also work (vim-style)"
echo "   - d - should NOT toggle dirs (global key disabled)"
echo "   - f - should NOT toggle files (global key disabled)"
echo "   - a - should NOT toggle aggregate (global key disabled)"
echo "   - s - should NOT cycle sort (global key disabled)"
echo "   - q - should close dialog"
echo "   - Esc - should close dialog"
echo ""
echo "🎯 Expected behavior:"
echo "- Arrow keys should move the '>' indicator"
echo "- Global keys (d,f,a,s) should be ignored in FileDialog"
echo "- Only FileDialog-specific keys should work"
echo ""
echo "📝 If arrow keys still don't work:"
echo "- Check if 'q' works (confirms keybindings are active)"
echo "- Check if global keys are blocked (d,f,a,s should do nothing)"
echo "- Try hjkl instead of arrow keys"
echo "- Check terminal type and key support" 