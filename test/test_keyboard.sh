#!/bin/bash

echo "🔧 Navigation Keys Test - watch-fs"
echo "====================================="
echo ""
echo "📋 Test instructions:"
echo "1. The application will start in 3 seconds"
echo "2. You should see events for the created files"
echo "3. Test the following keys:"
echo ""
echo "   Navigation:"
echo "   - ↑/↓/←/→ (arrows)"
echo "   - h/j/k/l (vim-style)"
echo "   - Page Up/Page Down"
echo "   - Home/End"
echo "   - g/G (beginning/end)"
echo ""
echo "   Others:"
echo "   - f : Toggle files"
echo "   - d : Toggle directories"
echo "   - a : Toggle aggregation"
echo "   - s : Change sorting"
echo "   - q : Quit"
echo ""
echo "⏳ Starting in 3 seconds..."
sleep 3

echo "🚀 Launching watch-fs..."
./bin/watch-fs -path . 