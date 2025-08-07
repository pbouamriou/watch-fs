#!/bin/bash

echo "=== Testing Folder Manager Currently Watching Debug ==="

# Compile
go build -o watch-fs ./cmd/watch-fs

echo "Starting watch-fs with debug output..."
echo "Will run for 5 seconds, then we'll simulate Ctrl+F press"

# Start the application in background and send Ctrl+F after a short delay
(
  sleep 2
  echo "Sending Ctrl+F to open folder manager..."
  # Unfortunately, we can't easily send keystrokes to TUI programmatically
  # We'll need to do this manually
) &

# Run the application with a timeout and capture output
timeout 5 ./watch-fs -path /Users/philippebouamriou/Documents 2>&1 | tee debug_output.log

echo ""
echo "Debug output captured. Check debug_output.log for details."
echo "Please manually test by running:"
echo "./watch-fs -path /Users/philippebouamriou/Documents"
echo "Then press Ctrl+F to see debug output."