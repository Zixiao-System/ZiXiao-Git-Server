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

# Check for Node.js and npm
if ! command -v node &> /dev/null; then
    echo "WARNING: Node.js is not installed"
    echo "Installing Node.js..."
    if [ "$OS" = "Darwin" ]; then
        if command -v brew &> /dev/null; then
            brew install node
        else
            echo "ERROR: Homebrew is not installed. Please install Node.js from https://nodejs.org/"
            exit 1
        fi
    else
        echo "Please install Node.js 18+ from https://nodejs.org/"
        exit 1
    fi
fi

NODE_VERSION=$(node --version)
echo "✓ Node.js found: $NODE_VERSION"

if ! command -v npm &> /dev/null; then
    echo "ERROR: npm is not installed"
    exit 1
fi

NPM_VERSION=$(npm --version)
echo "✓ npm found: $NPM_VERSION"

# Check for Nginx
if ! command -v nginx &> /dev/null; then
    echo "WARNING: Nginx is not installed"
    echo "Installing Nginx..."
    if [ "$OS" = "Darwin" ]; then
        if command -v brew &> /dev/null; then
            brew install nginx
        else
            echo "ERROR: Homebrew is not installed"
            echo "Please install Homebrew from https://brew.sh/"
            exit 1
        fi
    else
        # Linux
        if command -v apt-get &> /dev/null; then
            sudo apt-get update
            sudo apt-get install -y nginx
        elif command -v yum &> /dev/null; then
            sudo yum install -y nginx
        else
            echo "ERROR: Package manager not found. Please install Nginx manually."
            exit 1
        fi
    fi
fi

NGINX_VERSION=$(nginx -v 2>&1 | awk -F'/' '{print $2}')
echo "✓ Nginx found: $NGINX_VERSION"

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

# Install frontend dependencies
echo ""
echo "Installing frontend dependencies..."
cd frontend
npm install
echo "✓ Frontend dependencies installed"
cd ..

# Create directories
echo ""
echo "Creating directories..."
mkdir -p data/repositories
mkdir -p logs
mkdir -p bin
mkdir -p git-core/lib
mkdir -p web/dist
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

# Configure Nginx
echo ""
echo "Configuring Nginx..."

PROJECT_ROOT=$(pwd)

# Create local Nginx config with project path
sed "s|__PROJECT_ROOT__|$PROJECT_ROOT|g" configs/nginx-local.conf > configs/nginx-local-generated.conf

if [ "$OS" = "Darwin" ]; then
    # macOS Homebrew Nginx
    NGINX_CONF_DIR="/opt/homebrew/etc/nginx"
    if [ ! -d "$NGINX_CONF_DIR" ]; then
        NGINX_CONF_DIR="/usr/local/etc/nginx"
    fi

    if [ -d "$NGINX_CONF_DIR/servers" ]; then
        echo "Creating Nginx configuration symlink..."
        ln -sf "$PROJECT_ROOT/configs/nginx-local-generated.conf" "$NGINX_CONF_DIR/servers/zixiao-git-server.conf"
        echo "✓ Nginx configuration linked"
    else
        echo "WARNING: Nginx servers directory not found"
        echo "Please manually copy configs/nginx-local-generated.conf to your Nginx configuration"
    fi
else
    # Linux
    NGINX_CONF_DIR="/etc/nginx"
    if [ -d "$NGINX_CONF_DIR/sites-available" ]; then
        echo "Creating Nginx configuration..."
        sudo cp "$PROJECT_ROOT/configs/nginx-local-generated.conf" "$NGINX_CONF_DIR/sites-available/zixiao-git-server.conf"

        if [ ! -L "$NGINX_CONF_DIR/sites-enabled/zixiao-git-server.conf" ]; then
            sudo ln -s "$NGINX_CONF_DIR/sites-available/zixiao-git-server.conf" "$NGINX_CONF_DIR/sites-enabled/zixiao-git-server.conf"
        fi

        echo "✓ Nginx configuration installed"

        # Test Nginx config
        if sudo nginx -t &> /dev/null; then
            echo "✓ Nginx configuration is valid"
        else
            echo "WARNING: Nginx configuration test failed"
        fi
    else
        echo "WARNING: Nginx sites-available directory not found"
        echo "Please manually configure Nginx using configs/nginx-local-generated.conf"
    fi
fi

echo ""
echo "======================================"
echo "Installation completed successfully!"
echo "======================================"
echo ""
echo "Next steps:"
echo "  1. Review and update configs/server.yaml"
echo "  2. Build the frontend: cd frontend && npm run build"
echo "  3. Build the backend: make build"
echo "  4. Start the backend server: make run"
echo "  5. Start/Restart Nginx:"
if [ "$OS" = "Darwin" ]; then
    echo "     brew services start nginx"
    echo "     (or: brew services restart nginx)"
else
    echo "     sudo systemctl start nginx"
    echo "     (or: sudo systemctl restart nginx)"
fi
echo ""
echo "Access the application at: http://localhost"
echo ""
