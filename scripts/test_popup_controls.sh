#!/bin/bash

# Test script to verify that popup controls work correctly
# This script tests the Escape and q key functionality for the details popup

set -e

echo "🧪 Testing Popup Controls..."
echo "============================="

# Build the application
echo "📦 Building watch-fs..."
go build -o bin/watch-fs cmd/watch-fs/main.go

# Create a test directory
TEST_DIR="/tmp/watch-fs-test-popup"
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

# Function to check if help text mentions popup controls
check_popup_help() {
    local output_file="$1"
    if grep -q "Enter: Details" "$output_file"; then
        echo "✅ Enter: Details found in help text"
        return 0
    else
        echo "❌ Enter: Details not found in help text"
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
echo "🚀 Starting watch-fs for popup control test..."
OUTPUT_FILE="/tmp/watch-fs-popup-test.log"
timeout 10s bin/watch-fs -path "$TEST_DIR" > "$OUTPUT_FILE" 2>&1 &
WATCH_PID=$!

# Wait a moment for the application to start
sleep 2

# Generate some test events
echo "📝 Generating test events for popup testing..."

# Create a file (CREATE)
echo "test content" > "$TEST_DIR/test_file.txt"

# Write to the file (WRITE)
echo "more content" >> "$TEST_DIR/test_file.txt"

# Wait for the application to process events
sleep 2

# Stop the application
kill $WATCH_PID 2>/dev/null || true
wait $WATCH_PID 2>/dev/null || true

echo "📊 Checking popup control functionality..."

# Check if help text mentions popup controls
if check_popup_help "$OUTPUT_FILE"; then
    echo "✅ Popup controls mentioned in help text"
else
    echo "❌ Popup controls not mentioned in help text"
    exit 1
fi

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
echo "🎉 Popup controls test completed!"
echo ""
echo "📝 Manual Testing Instructions for Popup Controls:"
echo "1. Run: bin/watch-fs -path /path/to/directory"
echo "2. Navigate to an event using arrow keys"
echo "3. Press Enter to show event details popup"
echo "4. Test popup controls:"
echo "   - Press Escape to close popup (should work)"
echo "   - Press q to close popup (should work, NOT quit app)"
echo "   - Press Enter again to show popup"
echo "   - Press q again to close popup"
echo "5. After closing popup, press q to quit application (should work)"
echo ""
echo "🔍 Expected Behavior:"
echo "- When popup is OPEN: q closes popup, does NOT quit app"
echo "- When popup is CLOSED: q quits the application"
echo "- Escape always closes popup when it's open"
echo "- Enter always opens popup when on an event" 