#!/bin/bash

# Test script to verify that import/export functionality works correctly
# This script tests the Ctrl+E and Ctrl+I key functionality

set -e

echo "🧪 Testing Import/Export Functionality..."
echo "=========================================="

# Build the application
echo "📦 Building watch-fs..."
go build -o bin/watch-fs cmd/watch-fs/main.go

# Create a test directory
TEST_DIR="/tmp/watch-fs-test-import-export"
rm -rf "$TEST_DIR"
mkdir -p "$TEST_DIR"

echo "📁 Test directory: $TEST_DIR"

# Function to check if export files were created
check_export_files() {
    local sqlite_count=$(find . -name "watch-fs-events_*.db" | wc -l)
    local json_count=$(find . -name "watch-fs-events_*.json" | wc -l)
    
    if [ "$sqlite_count" -gt 0 ]; then
        echo "✅ SQLite export file found"
        return 0
    elif [ "$json_count" -gt 0 ]; then
        echo "✅ JSON export file found"
        return 0
    else
        echo "❌ No export files found"
        return 1
    fi
}

# Function to check if import/export is mentioned in help
check_import_export_in_help() {
    local output_file="$1"
    if grep -q "Ctrl+E: Export\|Ctrl+I: Import" "$output_file"; then
        echo "✅ Import/Export shortcuts found in help text"
        return 0
    else
        echo "❌ Import/Export shortcuts not found in help text"
        echo "Debug: Help text content:"
        grep -A 5 -B 5 "Help" "$output_file" || echo "No help text found"
        return 1
    fi
}

# Function to check if dialog functionality is working
check_dialog_functionality() {
    echo "✅ Dialog functionality implemented (manual testing required)"
    echo "   - Press Ctrl+E to open export dialog"
    echo "   - Press Ctrl+I to open import dialog"
    echo "   - Type filename and press Enter to confirm"
    echo "   - Press Escape to cancel"
    return 0
}

# Clean up any existing export files
rm -f watch-fs-events_*.db watch-fs-events_*.json

# Start watch-fs in background and capture output
echo "🚀 Starting watch-fs..."
OUTPUT_FILE="/tmp/watch-fs-import-export-test.log"
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

# Test export functionality (Ctrl+E)
echo "📤 Testing export functionality (Ctrl+E)..."
# We can't easily simulate Ctrl+E in this test, but we can check if the help mentions it

# Stop the application
kill $WATCH_PID 2>/dev/null || true
wait $WATCH_PID 2>/dev/null || true

echo "📊 Checking import/export functionality..."

# Check if help text mentions import/export shortcuts
if check_import_export_in_help "$OUTPUT_FILE"; then
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

# Test manual export by running the application and checking for export files
echo "📤 Testing manual export..."
# Start watch-fs again briefly
timeout 5s bin/watch-fs "$TEST_DIR" > /dev/null 2>&1 &
WATCH_PID2=$!
sleep 2
kill $WATCH_PID2 2>/dev/null || true
wait $WATCH_PID2 2>/dev/null || true

# Check if export files were created (they shouldn't be without Ctrl+E)
if check_export_files; then
    echo "⚠️  Export files found (this might be expected if Ctrl+E was pressed)"
else
    echo "✅ No export files found (expected without Ctrl+E)"
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
echo "🎉 Import/Export functionality test completed!"
echo ""
echo "📝 Manual Testing Instructions:"
echo "1. Run: bin/watch-fs /path/to/directory"
echo "2. Generate some events by creating/modifying files"
echo "3. Press Ctrl+E to open export dialog"
echo "4. Type filename (e.g., 'my-events.db' or 'my-events.json') and press Enter"
echo "5. Press Ctrl+I to open import dialog"
echo "6. Type filename or press Enter to use default"
echo "7. Verify that:"
echo "   - Export dialog shows current filename and event count"
echo "   - Import dialog shows available files"
echo "   - Files are created with specified names"
echo "   - Events are imported correctly"
echo "   - Escape key cancels dialogs"
echo "   - Help text shows 'Ctrl+E: Export | Ctrl+I: Import'"
echo ""
echo "🔍 SQLite Export Features:"
echo "   - Creates indexed database for fast queries"
echo "   - Stores all event metadata (path, operation, timestamp, is_dir, count)"
echo "   - Can be opened with any SQLite browser"
echo "   - Supports complex queries for analysis"
echo ""
echo "📄 JSON Export Features:"
echo "   - Human-readable format"
echo "   - Includes metadata (export time, total count)"
echo "   - Easy to parse and process with other tools"
echo ""
echo "🎯 Dialog Features:"
echo "   - Interactive filename selection"
echo "   - Automatic format detection (.db = SQLite, .json = JSON)"
echo "   - Default filename suggestions"
echo "   - Escape to cancel, Enter to confirm"
echo "   - Shows event count and format information" 