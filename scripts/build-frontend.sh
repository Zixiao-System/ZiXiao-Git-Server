#!/bin/bash

# Frontend build script for Linux/macOS

set -e

echo "======================================"
echo "Building Frontend"
echo "======================================"
echo ""

# Check if node_modules exists
if [ ! -d "frontend/node_modules" ]; then
    echo "Installing dependencies..."
    cd frontend
    npm install
    cd ..
fi

# Build frontend
echo "Building Vue application..."
cd frontend
npm run build
cd ..

echo ""
echo "✓ Frontend build completed successfully!"
echo "Built files are in: web/dist/"
echo ""
