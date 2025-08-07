#!/bin/bash

echo "=== Testing Fixed Folder Count Display ===="
echo ""
echo "Testing that each folder shows its specific file count, not the total"
echo ""
echo "Instructions:"
echo "1. Press Ctrl+F to open folder manager"
echo "2. Check the 'Currently Watching' panel (left side)"
echo "3. Each folder should show its own file count in parentheses"
echo "4. The counts should be DIFFERENT for each folder"
echo "5. Documents and Downloads should NOT show the same number"
echo ""
echo "Expected behavior:"
echo "- Documents (XXXX) - where XXXX is specific to Documents folder"
echo "- Downloads (YYYY) - where YYYY is specific to Downloads folder"
echo "- The numbers XXXX and YYYY should be different"
echo ""
echo "Previously broken behavior:"
echo "- Documents (9374) and Downloads (9374) - same total for both"
echo ""
echo "Press Enter to start..."
read

# Start with both Documents and Downloads to test different counts
./watch-fs -path /Users/philippebouamriou/Documents