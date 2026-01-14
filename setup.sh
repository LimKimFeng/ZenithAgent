#!/bin/bash

echo "========================================="
echo "  ZenithAgent Environment Setup"
echo "========================================="
echo ""

# Check if running as root
if [ "$EUID" -eq 0 ]; then 
    echo "⚠️  Please do not run this script as root"
    exit 1
fi

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check Node.js
echo "📦 Checking Node.js..."
if command_exists node; then
    NODE_VERSION=$(node --version)
    echo "✓ Node.js installed: $NODE_VERSION"
    if [[ "${NODE_VERSION:1:2}" -lt 18 ]]; then
        echo "⚠️  Warning: Node.js 18+ recommended (you have $NODE_VERSION)"
    fi
else
    echo "❌ Node.js not found!"
    echo "Install with: curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash - && sudo apt install -y nodejs"
    exit 1
fi

# Check npm
echo "📦 Checking npm..."
if command_exists npm; then
    NPM_VERSION=$(npm --version)
    echo "✓ npm installed: $NPM_VERSION"
else
    echo "❌ npm not found!"
    exit 1
fi

# Check Go
echo "📦 Checking Go..."
if command_exists go; then
    GO_VERSION=$(go version)
    echo "✓ Go installed: $GO_VERSION"
else
    echo "❌ Go not found!"
    echo "Install with: wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz"
    echo "              sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz"
    echo "              echo 'export PATH=\$PATH:/usr/local/go/bin' >> ~/.bashrc"
    exit 1
fi

# Check Tor (optional but recommended)
echo "📦 Checking Tor..."
if command_exists tor; then
    echo "✓ Tor installed"
else
    echo "⚠️  Tor not found (optional for IP rotation)"
    echo "Install with: sudo apt install -y tor"
fi

echo ""
echo "========================================="
echo "  Installing Dependencies"
echo "========================================="
echo ""

# Install Go dependencies
echo "📥 Installing Go modules..."
go mod download
if [ $? -eq 0 ]; then
    echo "✓ Go modules installed"
else
    echo "❌ Failed to install Go modules"
    exit 1
fi

# Install React dashboard dependencies
if [ -d "dashboard" ]; then
    echo "📥 Installing React dashboard dependencies..."
    cd dashboard
    npm install
    if [ $? -eq 0 ]; then
        echo "✓ React dependencies installed"
    else
        echo "❌ Failed to install React dependencies"
        exit 1
    fi
    cd ..
else
    echo "⚠️  Dashboard directory not found"
fi

echo ""
echo "========================================="
echo "  Building Application"
echo "========================================="
echo ""

# Build React dashboard
if [ -d "dashboard" ]; then
    echo "🔨 Building React dashboard..."
    cd dashboard
    npm run build
    if [ $? -eq 0 ]; then
        echo "✓ React dashboard built successfully"
    else
        echo "❌ React build failed"
        exit 1
    fi
    cd ..
fi

# Build Go binary
echo "🔨 Building Go binary..."
go build -ldflags="-s -w" -o sys-log-runtime cmd/agent/main.go
if [ $? -eq 0 ]; then
    echo "✓ Go binary built successfully"
else
    echo "❌ Go build failed"
    exit 1
fi

# Make run.sh executable
chmod +x run.sh
chmod +x start.sh

echo ""
echo "========================================="
echo "  ✅ Setup Complete!"
echo "========================================="
echo ""
echo "Next steps:"
echo "  1. Run './start.sh' to launch ZenithAgent"
echo "  2. Access dashboard at http://localhost:8080"
echo "  3. For VPS deployment, see DEPLOYMENT.md"
echo ""
echo "VPS Info:"
echo "  IP: 57.20.32.131"
echo "  SSH Port: 14034"
echo "  Dashboard will be at: http://zenith.linjinfeng.site"
echo ""
