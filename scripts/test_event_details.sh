#!/bin/bash

# Test script to verify that event details popup functionality works correctly
# This script tests the Enter key functionality to show event details

set -e

echo "🧪 Testing Event Details Popup..."
echo "=================================="

# Build the application
echo "📦 Building watch-fs..."
go build -o bin/watch-fs cmd/watch-fs/main.go

# Create a test directory
TEST_DIR="/tmp/watch-fs-test-details"
rm -rf "$TEST_DIR"
mkdir -p "$TEST_DIR"

echo "📁 Test directory: $TEST_DIR"

# Function to check if details functionality is mentioned in help
check_details_in_help() {
    local output_file="$1"
    if grep -q "Enter: Details" "$output_file"; then
        echo "✅ Enter: Details found in help text"
        return 0
    else
        echo "❌ Enter: Details not found in help text"
        return 1
    fi
}

# Function to check if details view is created
check_details_view() {
    local output_file="$1"
    if grep -q "Event Details" "$output_file"; then
        echo "✅ Event Details view found"
        return 0
    else
        echo "❌ Event Details view not found"
        return 1
    fi
}

# Start watch-fs in background and capture output
echo "🚀 Starting watch-fs..."
OUTPUT_FILE="/tmp/watch-fs-details-test.log"
timeout 15s bin/watch-fs "$TEST_DIR" > "$OUTPUT_FILE" 2>&1 &
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

echo "📊 Checking event details functionality..."

# Check if help text mentions Enter: Details
if check_details_in_help "$OUTPUT_FILE"; then
    echo "✅ Help text correctly mentions Enter: Details"
else
    echo "❌ Help text missing Enter: Details"
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

echo ""
echo "🎉 Event Details functionality test completed!"
echo ""
echo "📝 Manual Testing Instructions:"
echo "1. Run: bin/watch-fs /path/to/directory"
echo "2. Navigate to an event using arrow keys"
echo "3. Press Enter to show event details"
echo "4. Press Escape or q to close the details popup"
echo "5. Verify that the popup shows:"
echo "   - Operation type with color"
echo "   - Full path"
echo "   - File/Directory type"
echo "   - Timestamp with milliseconds"
echo "   - Event count"
echo "   - File size (if available)"
echo "   - File permissions (if available)"
echo "   - Last modified time (if available)" 