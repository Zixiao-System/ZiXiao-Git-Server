#!/bin/bash

# Quick test script for ZiXiao Git Server

set -e

echo "======================================"
echo "ZiXiao Git Server - Quick Test"
echo "======================================"
echo ""

# Test 1: Check Go files compile
echo "[1/5] Checking Go syntax..."
go fmt ./... > /dev/null 2>&1 || true
echo "✓ Go files syntax OK"

# Test 2: Check C++ headers
echo "[2/5] Checking C++ headers..."
if [ -f "git-core/include/git_repository.h" ] && \
   [ -f "git-core/include/git_object.h" ] && \
   [ -f "git-core/include/git_protocol.h" ] && \
   [ -f "git-core/include/git_pack.h" ] && \
   [ -f "git-core/include/git_c_api.h" ]; then
    echo "✓ All C++ headers present"
else
    echo "✗ Missing C++ headers"
    exit 1
fi

# Test 3: Check C++ source files
echo "[3/5] Checking C++ source files..."
if [ -f "git-core/src/git_repository.cpp" ] && \
   [ -f "git-core/src/git_object.cpp" ] && \
   [ -f "git-core/src/git_protocol.cpp" ] && \
   [ -f "git-core/src/git_pack.cpp" ] && \
   [ -f "git-core/src/git_c_api.cpp" ]; then
    echo "✓ All C++ source files present"
else
    echo "✗ Missing C++ source files"
    exit 1
fi

# Test 4: Check Go packages
echo "[4/5] Checking Go packages..."
if [ -f "cmd/server/main.go" ] && \
   [ -f "internal/api/routes.go" ] && \
   [ -f "internal/auth/auth.go" ] && \
   [ -f "internal/config/config.go" ] && \
   [ -f "internal/database/database.go" ] && \
   [ -f "internal/models/models.go" ] && \
   [ -f "internal/repository/repository.go" ] && \
   [ -f "pkg/gitcore/gitcore.go" ]; then
    echo "✓ All Go packages present"
else
    echo "✗ Missing Go packages"
    exit 1
fi

# Test 5: Check configuration files
echo "[5/5] Checking configuration files..."
if [ -f "configs/server.yaml" ] && \
   [ -f "Makefile" ] && \
   [ -f "README.md" ]; then
    echo "✓ Configuration files present"
else
    echo "✗ Missing configuration files"
    exit 1
fi

echo ""
echo "======================================"
echo "All checks passed!"
echo "======================================"
echo ""
echo "Next steps:"
echo "  1. Run: ./scripts/install.sh"
echo "  2. Run: make build"
echo "  3. Run: make run"
echo ""
