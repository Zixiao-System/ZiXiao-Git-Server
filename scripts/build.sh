#!/bin/bash

# Build script for ZiXiao Git Server
# This script builds both the C++ Git core library and the Go server

set -e

echo "======================================"
echo "ZiXiao Git Server Build Script"
echo "======================================"
echo ""

# Detect OS
OS="$(uname -s)"
echo "Detected OS: $OS"

# Create necessary directories
echo "Creating directories..."
mkdir -p git-core/lib
mkdir -p bin
mkdir -p data/repositories
mkdir -p logs

# Build C++ library
echo ""
echo "Building C++ Git core library..."

CXX=g++
CXXFLAGS="-std=c++17 -Wall -O2 -fPIC"
INCLUDES="-I./git-core/include"
LDFLAGS="-shared -lssl -lcrypto -lz"
LIBNAME="libgitcore.so"

# macOS specific settings
if [ "$OS" = "Darwin" ]; then
    LIBNAME="libgitcore.dylib"
    LDFLAGS="-dynamiclib -lssl -lcrypto -lz"

    # Check for Homebrew OpenSSL
    if [ -d "/opt/homebrew/opt/openssl" ]; then
        CXXFLAGS="$CXXFLAGS -I/opt/homebrew/opt/openssl/include"
        LDFLAGS="$LDFLAGS -L/opt/homebrew/opt/openssl/lib"
    elif [ -d "/usr/local/opt/openssl" ]; then
        CXXFLAGS="$CXXFLAGS -I/usr/local/opt/openssl/include"
        LDFLAGS="$LDFLAGS -L/usr/local/opt/openssl/lib"
    fi
fi

# Compile C++ source files
echo "Compiling C++ source files..."
$CXX $CXXFLAGS $INCLUDES -c git-core/src/git_repository.cpp -o git-core/src/git_repository.o
$CXX $CXXFLAGS $INCLUDES -c git-core/src/git_object.cpp -o git-core/src/git_object.o
$CXX $CXXFLAGS $INCLUDES -c git-core/src/git_protocol.cpp -o git-core/src/git_protocol.o
$CXX $CXXFLAGS $INCLUDES -c git-core/src/git_pack.cpp -o git-core/src/git_pack.o
$CXX $CXXFLAGS $INCLUDES -c git-core/src/git_c_api.cpp -o git-core/src/git_c_api.o

# Link shared library
echo "Linking shared library..."
$CXX $LDFLAGS -o git-core/lib/$LIBNAME \
    git-core/src/git_repository.o \
    git-core/src/git_object.o \
    git-core/src/git_protocol.o \
    git-core/src/git_pack.o \
    git-core/src/git_c_api.o

echo "C++ library built: git-core/lib/$LIBNAME"

# Build Go server
echo ""
echo "Building Go server..."
CGO_ENABLED=1 go build -o bin/zixiao-git-server ./cmd/server

echo ""
echo "======================================"
echo "Build completed successfully!"
echo "======================================"
echo ""
echo "To run the server:"
echo "  ./bin/zixiao-git-server -config ./configs/server.yaml"
echo ""
