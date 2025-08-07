#!/bin/bash

echo "=== Testing Fixed Currently Watching Display ==="
echo ""
echo "Starting watch-fs with Documents folder..."
echo "The left panel should now show 'Documents' in Currently Watching"
echo ""
echo "Instructions:"
echo "1. Press Ctrl+F to open folder manager"
echo "2. Check the LEFT panel for 'Currently Watching' content"
echo "3. Try using arrow keys in left panel to navigate (if multiple folders)"
echo "4. Try 'r' to remove a folder from watching"
echo "5. Press Esc to close"
echo ""
echo "Press Enter to start..."
read

./watch-fs -path /Users/philippebouamriou/Documents