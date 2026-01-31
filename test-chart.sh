#!/bin/bash
# Quick test script to preview the chart feature

echo "Building nes-outage-status-checker..."
if command -v go &> /dev/null; then
    go build -o nes-outage-status-checker .
    if [ $? -eq 0 ]; then
        echo "✓ Build successful!"
        echo ""
        echo "To test the chart feature:"
        echo "1. Run: ./nes-outage-status-checker <event-id>"
        echo "2. Wait for a few auto-refreshes (every 30 seconds)"
        echo "3. Press 'c' to toggle to chart view"
        echo "4. Press 'c' again to return to detail view"
        echo ""
        echo "To find an event ID, visit: https://www.nespower.com/outages/"
    else
        echo "✗ Build failed"
        exit 1
    fi
else
    echo "✗ Go is not installed or not in PATH"
    echo "Install Go: brew install go"
    exit 1
fi
