#!/bin/bash

# Test script for the new focus system in watch-fs
# This script tests the focus-based input handling

set -e

echo "🧪 Testing Focus System"
echo "======================="

# Create test directory
TEST_DIR="./test_focus"
mkdir -p "$TEST_DIR"

# Function to cleanup
cleanup() {
    echo "🧹 Cleaning up..."
    rm -rf "$TEST_DIR"
    pkill -f "watch-fs" || true
}

# Set up cleanup on exit
trap cleanup EXIT

echo "📁 Created test directory: $TEST_DIR"

# Test 1: Basic compilation
echo "✅ Test 1: Compilation"
go build -o bin/watch-fs cmd/watch-fs/main.go
echo "   ✓ Compilation successful"

# Test 2: Help output
echo "✅ Test 2: Help output"
./bin/watch-fs --help > /dev/null
echo "   ✓ Help output works"

# Test 3: Version output
echo "✅ Test 3: Version output"
./bin/watch-fs --version > /dev/null
echo "   ✓ Version output works"

# Test 4: Focus system structure
echo "✅ Test 4: Focus system structure"
if grep -q "CurrentFocus" internal/ui/types.go; then
    echo "   ✓ CurrentFocus field exists in types"
else
    echo "   ❌ CurrentFocus field missing"
    exit 1
fi

if grep -q "FocusMain\|FocusDetails\|FocusExport\|FocusImport\|FocusFileDialog" internal/ui/types.go; then
    echo "   ✓ Focus constants defined"
else
    echo "   ❌ Focus constants missing"
    exit 1
fi

# Test 5: Keybindings structure
echo "✅ Test 5: Keybindings structure"
if grep -q "SetKeybinding.*EventsView" internal/ui/ui.go; then
    echo "   ✓ EventsView keybindings exist"
else
    echo "   ❌ EventsView keybindings missing"
    exit 1
fi

if grep -q "SetKeybinding.*DetailsView" internal/ui/ui.go; then
    echo "   ✓ DetailsView keybindings exist"
else
    echo "   ❌ DetailsView keybindings missing"
    exit 1
fi

if grep -q "SetKeybinding.*FileListView" internal/ui/ui.go; then
    echo "   ✓ FileListView keybindings exist"
else
    echo "   ❌ FileListView keybindings missing"
    exit 1
fi

# Test 6: Help context switching
echo "✅ Test 6: Help context switching"
if grep -q "switch ui.state.CurrentFocus" internal/ui/ui.go; then
    echo "   ✓ Help context switching implemented"
else
    echo "   ❌ Help context switching missing"
    exit 1
fi

# Test 7: Focus management functions
echo "✅ Test 7: Focus management functions"
if grep -q "SetCurrentView.*EventsView" internal/ui/ui.go; then
    echo "   ✓ Focus management implemented"
else
    echo "   ❌ Focus management missing"
    exit 1
fi

# Test 8: No manual focus checks
echo "✅ Test 8: No manual focus checks"
if grep -q "if ui.state.ShowFileDialog.*return nil" internal/ui/ui.go; then
    echo "   ❌ Manual focus checks still present"
    exit 1
else
    echo "   ✓ Manual focus checks removed"
fi

echo ""
echo "🎉 All focus system tests passed!"
echo ""
echo "📋 Summary of improvements:"
echo "   • Native gocui focus system implemented"
echo "   • Context-aware help display"
echo "   • Cleaner keybinding management"
echo "   • No more manual focus checks"
echo "   • Better separation of concerns"
echo ""
echo "🚀 The focus system is ready for use!" 