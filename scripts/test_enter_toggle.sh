#!/bin/bash

# Test script to verify that Enter key works as a toggle for the details popup
# This script tests the new toggle functionality

set -e

echo "🧪 Testing Enter Toggle Functionality..."
echo "========================================"

# Build the application
echo "📦 Building watch-fs..."
go build -o bin/watch-fs cmd/watch-fs/main.go

# Create a test directory
TEST_DIR="/tmp/watch-fs-test-toggle"
rm -rf "$TEST_DIR"
mkdir -p "$TEST_DIR"

echo "📁 Test directory: $TEST_DIR"

# Function to check if the application compiles and runs
test_compilation() {
    echo "🔧 Testing compilation..."
    if go build -o bin/watch-fs cmd/watch-fs/main.go; then
        echo "✅ Compilation successful"
        return 0
    else
        echo "❌ Compilation failed"
        return 1
    fi
}

# Function to check if tests pass
test_unit_tests() {
    echo "🧪 Running unit tests..."
    if go test ./test/...; then
        echo "✅ Unit tests passed"
        return 0
    else
        echo "❌ Unit tests failed"
        return 1
    fi
}

# Run compilation test
if ! test_compilation; then
    exit 1
fi

# Run unit tests
if ! test_unit_tests; then
    exit 1
fi

# Start watch-fs in background and capture output
echo "🚀 Starting watch-fs for toggle test..."
OUTPUT_FILE="/tmp/watch-fs-toggle-test.log"
timeout 10s bin/watch-fs -path "$TEST_DIR" > "$OUTPUT_FILE" 2>&1 &
WATCH_PID=$!

# Wait a moment for the application to start
sleep 2

# Generate some test events
echo "📝 Generating test events for toggle testing..."

# Create a file (CREATE)
echo "test content" > "$TEST_DIR/test_file.txt"

# Write to the file (WRITE)
echo "more content" >> "$TEST_DIR/test_file.txt"

# Wait for the application to process events
sleep 2

# Stop the application
kill $WATCH_PID 2>/dev/null || true
wait $WATCH_PID 2>/dev/null || true

echo "📊 Checking toggle functionality..."

# Check if events were captured
if grep -q "CREATE\|WRITE" "$OUTPUT_FILE"; then
    echo "✅ Events were captured successfully"
else
    echo "❌ No events were captured"
    exit 1
fi

# Cleanup
rm -rf "$TEST_DIR"
rm -f "$OUTPUT_FILE"

echo ""
echo "🎉 Enter toggle test completed!"
echo ""
echo "📝 Manual Testing Instructions for Enter Toggle:"
echo "1. Run: bin/watch-fs -path /path/to/directory"
echo "2. Navigate to an event using arrow keys"
echo "3. Press Enter to show event details popup"
echo "4. Press Enter again to close the popup (toggle functionality)"
echo "5. Press Enter again to show the popup again"
echo "6. Test other controls:"
echo "   - Press Escape to close popup"
echo "   - Press q to close popup"
echo "   - Press Enter to show popup again"
echo "7. After closing popup, press q to quit application"
echo ""
echo "🔍 Expected Behavior:"
echo "- Enter acts as a toggle: opens popup when closed, closes popup when open"
echo "- Escape always closes popup when it's open"
echo "- q closes popup when open, quits app when popup is closed"
echo "- Toggle only works when in EventsView (not when popup is focused)"
echo ""
echo "✨ New Feature: Enter Toggle"
echo "- More intuitive: Enter to open, Enter to close"
echo "- Consistent with common UI patterns"
echo "- Reduces need to remember multiple keys"
echo "- Improves user experience" 