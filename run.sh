#!/bin/bash

# Build with ldflags to strip symbols (-s -w) for smaller binary and slight obfuscation
echo "Building ZenithAgent..."
# Ensure we have dependencies (if go is installed)
if command -v go &> /dev/null; then
    go mod tidy
    go build -ldflags="-s -w" -o sys-log-runtime cmd/agent/main.go
else
    echo "Error: 'go' command not found. Please install Go."
    exit 1
fi

if [ -f "./sys-log-runtime" ]; then
    echo "Build successful. Launching stealth agent..."
    ./sys-log-runtime
else
    echo "Build failed."
    exit 1
fi
