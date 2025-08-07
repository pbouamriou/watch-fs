#!/bin/bash

echo "=== Testing Multiple --path Flags Support ===="
echo ""
echo "Testing the new ability to use multiple --path flags in command line"
echo ""
echo "NEW FUNCTIONALITY:"
echo "- watch-fs --path /dir1 --path /dir2 --path /dir3"
echo "- Each --path can be used multiple times"
echo "- Replaces the need for comma-separated --paths flag"
echo ""

# Test 1: Single path (should work as before)
echo "TEST 1: Single --path flag"
echo "Command: ./watch-fs --path /Users/philippebouamriou/Documents --version"
./watch-fs --path /Users/philippebouamriou/Documents --version
echo ""

# Test 2: Multiple paths 
echo "TEST 2: Multiple --path flags"
echo "Command: ./watch-fs --path /Users/philippebouamriou/Documents --path /Users/philippebouamriou/Downloads --version"
./watch-fs --path /Users/philippebouamriou/Documents --path /Users/philippebouamriou/Downloads --version
echo ""

# Test 3: Show help to verify usage message
echo "TEST 3: Help message (should show new usage)"
echo "Command: ./watch-fs --help"
./watch-fs --help
echo ""

echo "VERIFICATION:"
echo "1. All commands above should work without errors"
echo "2. Help message should show the new usage patterns"
echo "3. Version should display correctly"
echo ""

echo "INTERACTIVE TEST:"
echo "Now testing with actual TUI... Press Ctrl+F to open folder manager"
echo "You should see BOTH Documents and Downloads in 'Currently Watching'"
echo ""
echo "Command: ./watch-fs --path /Users/philippebouamriou/Documents --path /Users/philippebouamriou/Downloads"
echo "Press Enter to start interactive test..."
read

./watch-fs --path /Users/philippebouamriou/Documents --path /Users/philippebouamriou/Downloads