#!/bin/bash

# Installation script for ZiXiao Git Server
# This script checks dependencies and sets up the environment

set -e

echo "======================================"
echo "ZiXiao Git Server Installation"
echo "======================================"
echo ""

# Detect OS
OS="$(uname -s)"
echo "Detected OS: $OS"

# Check for required tools
echo ""
echo "Checking dependencies..."

# Check for Go
if ! command -v go &> /dev/null; then
    echo "ERROR: Go is not installed"
    echo "Please install Go 1.21 or higher from https://golang.org/dl/"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}')
echo "✓ Go found: $GO_VERSION"

# Check for C++ compiler
if ! command -v g++ &> /dev/null; then
    echo "ERROR: g++ compiler is not installed"
    if [ "$OS" = "Darwin" ]; then
        echo "Please install Xcode Command Line Tools: xcode-select --install"
    else
        echo "Please install g++: apt-get install g++ (Ubuntu/Debian) or yum install gcc-c++ (RHEL/CentOS)"
    fi
    exit 1
fi

CXX_VERSION=$(g++ --version | head -n1)
echo "✓ C++ compiler found: $CXX_VERSION"

# Check for OpenSSL
if [ "$OS" = "Darwin" ]; then
    if [ ! -d "/opt/homebrew/opt/openssl" ] && [ ! -d "/usr/local/opt/openssl" ]; then
        echo "WARNING: OpenSSL not found via Homebrew"
        echo "Installing OpenSSL via Homebrew..."
        if command -v brew &> /dev/null; then
            brew install openssl
        else
            echo "ERROR: Homebrew is not installed"
            echo "Please install Homebrew from https://brew.sh/"
            exit 1
        fi
    else
        echo "✓ OpenSSL found"
    fi
else
    if ! ldconfig -p | grep -q libssl; then
        echo "WARNING: OpenSSL library not found"
        echo "Please install OpenSSL development libraries"
        exit 1
    else
        echo "✓ OpenSSL found"
    fi
fi

# Check for zlib
if [ "$OS" = "Darwin" ]; then
    echo "✓ zlib available (system)"
else
    if ! ldconfig -p | grep -q libz; then
        echo "WARNING: zlib library not found"
        echo "Please install zlib development libraries"
        exit 1
    else
        echo "✓ zlib found"
    fi
fi

# Install Go dependencies
echo ""
echo "Installing Go dependencies..."
go mod download
echo "✓ Go dependencies installed"

# Create directories
echo ""
echo "Creating directories..."
mkdir -p data/repositories
mkdir -p logs
mkdir -p bin
mkdir -p git-core/lib
echo "✓ Directories created"

# Create default config if it doesn't exist
if [ ! -f "configs/server.yaml" ]; then
    echo ""
    echo "Creating default configuration..."
    cat > configs/server.yaml << EOF
server:
  host: 0.0.0.0
  port: 8080
  mode: release

database:
  type: sqlite
  path: ./data/gitserver.db

git:
  repo_path: ./data/repositories
  max_repo_size: 1024
  max_file_size: 100

security:
  jwt_secret: $(openssl rand -hex 32)
  jwt_expiration: 24
  password_min: 8
  enable_ssh: false
  ssh_port: 2222
EOF
    echo "✓ Default configuration created: configs/server.yaml"
    echo "  IMPORTANT: Update the jwt_secret in production!"
fi

echo ""
echo "======================================"
echo "Installation completed successfully!"
echo "======================================"
echo ""
echo "Next steps:"
echo "  1. Review and update configs/server.yaml"
echo "  2. Build the project: make build"
echo "  3. Run the server: make run"
echo ""
