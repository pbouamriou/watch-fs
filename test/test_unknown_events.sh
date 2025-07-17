#!/bin/bash

# Test script to verify that UNKNOWN events are no longer appearing
# This script tests the fix for combined fsnotify operations

set -e

echo "🧪 Testing UNKNOWN events fix..."
echo "=================================="

# Build the application
echo "📦 Building watch-fs..."
go build -o bin/watch-fs cmd/watch-fs/main.go

# Create a test directory
TEST_DIR="/tmp/watch-fs-test-unknown"
rm -rf "$TEST_DIR"
mkdir -p "$TEST_DIR"

echo "📁 Test directory: $TEST_DIR"

# Function to check for UNKNOWN events
check_unknown_events() {
    local output_file="$1"
    if grep -q "UNKNOWN" "$output_file"; then
        echo "❌ UNKNOWN events found in output:"
        grep "UNKNOWN" "$output_file"
        return 1
    else
        echo "✅ No UNKNOWN events found"
        return 0
    fi
}

# Start watch-fs in background and capture output
echo "🚀 Starting watch-fs..."
OUTPUT_FILE="/tmp/watch-fs-unknown-test.log"
timeout 10s bin/watch-fs "$TEST_DIR" > "$OUTPUT_FILE" 2>&1 &
WATCH_PID=$!

# Wait a moment for the application to start
sleep 2

# Generate various file system events that might trigger combined operations
echo "📝 Generating test events..."

# Create a file (CREATE)
echo "test" > "$TEST_DIR/test1.txt"

# Write to the file (WRITE)
echo "more content" >> "$TEST_DIR/test1.txt"

# Change permissions (CHMOD)
chmod 755 "$TEST_DIR/test1.txt"

# Create a directory (CREATE)
mkdir "$TEST_DIR/testdir"

# Create a file in the directory (CREATE)
echo "test" > "$TEST_DIR/testdir/test2.txt"

# Remove the file (REMOVE)
rm "$TEST_DIR/test1.txt"

# Rename the directory (RENAME)
mv "$TEST_DIR/testdir" "$TEST_DIR/renamed_dir"

# Wait for the application to process events
sleep 3

# Stop the application
kill $WATCH_PID 2>/dev/null || true
wait $WATCH_PID 2>/dev/null || true

echo "📊 Checking for UNKNOWN events..."
if check_unknown_events "$OUTPUT_FILE"; then
    echo "✅ Test PASSED: No UNKNOWN events detected"
    echo ""
    echo "📋 Summary of events captured:"
    if [ -f "$OUTPUT_FILE" ]; then
        echo "Events found:"
        grep -E "(CREATE|WRITE|REMOVE|RENAME|CHMOD)" "$OUTPUT_FILE" | head -10
    fi
else
    echo "❌ Test FAILED: UNKNOWN events detected"
    exit 1
fi

# Cleanup
rm -rf "$TEST_DIR"
rm -f "$OUTPUT_FILE"

echo ""
echo "🎉 UNKNOWN events fix test completed successfully!" 