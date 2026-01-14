#!/bin/bash

echo "========================================="
echo "  🚀 Starting ZenithAgent"
echo "========================================="
echo ""

# Check if binary exists
if [ ! -f "./sys-log-runtime" ]; then
    echo "❌ Binary not found!"
    echo "Please run './setup.sh' first to build the application."
    exit 1
fi

# Check if React build exists
if [ ! -d "./dashboard/dist" ]; then
    echo "⚠️  React build not found. Building dashboard..."
    if [ -d "dashboard" ]; then
        cd dashboard
        npm run build
        if [ $? -ne 0 ]; then
            echo "❌ React build failed"
            echo "Continuing with embedded HTML fallback..."
        else
            echo "✓ React dashboard built"
        fi
        cd ..
    fi
fi

# Check for lock file from previous run
if [ -f ".zenith.lock" ]; then
    echo "⚠️  Found lock file from previous run"
    read -p "Remove lock file and continue? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        rm .zenith.lock
        echo "✓ Lock file removed"
    else
        echo "Exiting. Please manually remove .zenith.lock if needed."
        exit 1
    fi
fi

echo "Starting ZenithAgent..."
echo "Dashboard will be available at:"
echo "  - Local: http://localhost:8080"
echo "  - VPS: http://zenith.linjinfeng.site (after deployment)"
echo ""
echo "Press Ctrl+C to stop"
echo "========================================="
echo ""

# Run the application
./sys-log-runtime
