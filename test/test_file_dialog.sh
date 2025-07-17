#!/bin/bash

# Test script to verify that file dialog functionality works correctly
# This script tests the new file dialog interface for import/export

set -e

echo "🧪 Testing File Dialog Functionality..."
echo "========================================"

# Build the application
echo "📦 Building watch-fs..."
go build -o bin/watch-fs cmd/watch-fs/main.go

# Create a test directory
TEST_DIR="/tmp/watch-fs-test-file-dialog"
rm -rf "$TEST_DIR"
mkdir -p "$TEST_DIR"

echo "📁 Test directory: $TEST_DIR"

# Function to check if file dialog functionality is mentioned in help
check_file_dialog_in_help() {
    local output_file="$1"
    if grep -q "Ctrl+E.*Export\|Ctrl+I.*Import" "$output_file"; then
        echo "✅ Import/Export shortcuts found in help text"
        return 0
    else
        echo "❌ Import/Export shortcuts not found in help text"
        echo "Debug: Help text content:"
        grep -A 5 -B 5 "Help" "$output_file" || echo "No help text found"
        return 1
    fi
}

# Function to check if file dialog navigation is mentioned
check_file_dialog_navigation() {
    local output_file="$1"
    if grep -q "↑↓.*Navigate\|Enter.*Select\|Esc.*Cancel" "$output_file"; then
        echo "✅ File dialog navigation found in help text"
        return 0
    else
        echo "❌ File dialog navigation not found in help text"
        return 1
    fi
}

# Clean up any existing export files
rm -f watch-fs-events_*.db watch-fs-events_*.json

# Start watch-fs in background and capture output
echo "🚀 Starting watch-fs..."
OUTPUT_FILE="/tmp/watch-fs-file-dialog-test.log"
timeout 20s bin/watch-fs "$TEST_DIR" > "$OUTPUT_FILE" 2>&1 &
WATCH_PID=$!

# Wait a moment for the application to start
sleep 2

# Generate some test events
echo "📝 Generating test events..."

# Create a file (CREATE)
echo "test content" > "$TEST_DIR/test_file.txt"

# Write to the file (WRITE)
echo "more content" >> "$TEST_DIR/test_file.txt"

# Create a directory (CREATE)
mkdir "$TEST_DIR/test_dir"

# Create a file in the directory (CREATE)
echo "test" > "$TEST_DIR/test_dir/nested_file.txt"

# Wait for the application to process events
sleep 3

# Stop the application
kill $WATCH_PID 2>/dev/null || true
wait $WATCH_PID 2>/dev/null || true

echo "📊 Checking file dialog functionality..."

# Check if help text mentions import/export shortcuts
if check_file_dialog_in_help "$OUTPUT_FILE"; then
    echo "✅ Help text correctly mentions import/export shortcuts"
else
    echo "❌ Help text missing import/export shortcuts"
    exit 1
fi

# Check if events were captured
if grep -q "CREATE\|WRITE" "$OUTPUT_FILE"; then
    echo "✅ Events were captured successfully"
else
    echo "❌ No events were captured"
    exit 1
fi

echo ""
echo "📋 Summary of events captured:"
if [ -f "$OUTPUT_FILE" ]; then
    echo "Events found:"
    grep -E "(CREATE|WRITE|REMOVE|RENAME|CHMOD)" "$OUTPUT_FILE" | head -5
fi

# Cleanup
rm -rf "$TEST_DIR"
rm -f "$OUTPUT_FILE"
rm -f watch-fs-events_*.db watch-fs-events_*.json

echo ""
echo "🎉 File Dialog functionality test completed!"
echo ""
echo "📝 Manual Testing Instructions:"
echo "1. Run: bin/watch-fs /path/to/directory"
echo "2. Generate some events by creating/modifying files"
echo "3. Press Ctrl+E to open file dialog (save mode)"
echo "4. Navigate using arrow keys or hjkl"
echo "5. Press Enter to open directories or select files"
echo "6. Press Escape to cancel"
echo "7. Press Ctrl+I to open file dialog (open mode)"
echo "8. Browse to find export files"
echo "9. Select file and press Enter to import"
echo "10. Verify that:"
echo "    - File dialog shows current path"
echo "    - Files and directories are listed with icons"
echo "    - Navigation works with arrow keys and hjkl"
echo "    - Enter opens directories or selects files"
echo "    - Escape cancels the dialog"
echo "    - Help text shows navigation shortcuts when dialog is open"
echo ""
echo "🎯 File Dialog Features:"
echo "   - Modern file browser interface"
echo "   - Directory navigation with arrow keys"
echo "   - File listing with sizes and icons"
echo "   - Path display at top of dialog"
echo "   - Filter support for .db and .json files"
echo "   - Visual selection highlighting"
echo "   - Parent directory navigation (..)"
echo "   - Hidden file filtering"
echo "   - Sort: directories first, then files alphabetically" 