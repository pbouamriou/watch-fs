#!/bin/bash

# Test script for FileDialog save mode with filename editing
# This script tests the new filename editing feature in save mode

set -e

echo "🧪 Testing FileDialog save mode with filename editing..."

# Build the application
echo "📦 Building application..."
go build -o bin/watch-fs cmd/watch-fs/main.go

# Create test files
echo "📁 Creating test files..."
mkdir -p test_save
echo "test1" > test_save/file1.txt
echo "test2" > test_save/file2.txt
mkdir -p test_save/subdir
echo "subtest" > test_save/subdir/subfile.txt

echo "✅ Test setup completed!"
echo ""
echo "🔧 Testing steps:"
echo "1. Run: ./bin/watch-fs test_save"
echo "2. Press Ctrl+E to open export FileDialog (Save mode)"
echo "3. Test the new functionality:"
echo "   - Navigate with ↑↓/kj"
echo "   - Press 'e' to edit filename (should switch to filename input)"
echo "   - Type a custom filename (e.g., 'my-custom-export.db')"
echo "   - Press Enter to save with custom name"
echo "   - Press Esc to cancel editing and return to file list"
echo "4. Test file selection:"
echo "   - Select an existing file with Enter"
echo "   - Should overwrite the selected file"
echo ""
echo "🎯 Expected behavior:"
echo "- 'e' key should switch to filename editing mode"
echo "- Filename input should be editable"
echo "- Custom filenames should be saved correctly"
echo "- Enter should save the file with custom name"
echo "- Esc should cancel editing and return to file list"
echo "- File selection should still work normally"
echo ""
echo "📝 Features to test:"
echo "- Custom filename with .db extension"
echo "- Custom filename with .json extension"
echo "- Custom filename without extension (should add .db)"
echo "- Navigation between file list and filename input" 