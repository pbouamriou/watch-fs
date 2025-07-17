#!/bin/bash

echo "Debug test for 'e' key in FileDialog"
echo "This will show debug output when pressing 'e'"
echo ""

# Start watch-fs with debug output
echo "Starting watch-fs..."
./bin/watch-fs . 2>&1 | tee debug_output.log &
WATCH_FS_PID=$!

# Wait for startup
sleep 2

echo "Sending Ctrl+E to open export dialog..."
osascript -e 'tell application "System Events" to keystroke "e" using control down'

sleep 1

echo "Sending 'e' to test filename editing..."
osascript -e 'tell application "System Events" to keystroke "e"'

sleep 2

echo "Sending Ctrl+C to quit..."
osascript -e 'tell application "System Events" to keystroke "c" using control down'

wait $WATCH_FS_PID

echo ""
echo "Debug output:"
cat debug_output.log

echo ""
echo "Test completed. Check debug_output.log for debug messages." 