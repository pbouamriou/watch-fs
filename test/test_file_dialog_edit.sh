#!/bin/bash

# Test script for FileDialog edit functionality
# This script tests that pressing 'e' in the FileDialog (Save mode) switches to filename editing

echo "Testing FileDialog edit functionality..."
echo "This test will:"
echo "1. Start watch-fs"
echo "2. Press Ctrl+E to open export dialog"
echo "3. Press 'e' to switch to filename editing"
echo "4. Verify that the focus switches to filename input"
echo ""

# Start watch-fs in background
echo "Starting watch-fs..."
./bin/watch-fs . &
WATCH_FS_PID=$!

# Wait a moment for the application to start
sleep 2

echo "Sending Ctrl+E to open export dialog..."
# Send Ctrl+E to open export dialog
osascript -e 'tell application "System Events" to keystroke "e" using control down'

# Wait a moment for the dialog to appear
sleep 1

echo "Sending 'e' to switch to filename editing..."
# Send 'e' to switch to filename editing
osascript -e 'tell application "System Events" to keystroke "e"'

# Wait a moment to see the result
sleep 2

echo "Sending Ctrl+C to quit..."
# Send Ctrl+C to quit
osascript -e 'tell application "System Events" to keystroke "c" using control down'

# Wait for the process to terminate
wait $WATCH_FS_PID

echo ""
echo "Test completed!"
echo "If the test worked correctly, you should have seen:"
echo "1. The export dialog open when pressing Ctrl+E"
echo "2. The focus switch to filename input when pressing 'e'"
echo "3. The ability to type in the filename field" 