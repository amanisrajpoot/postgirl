#!/bin/bash

# Build script for Postgirl - Cross-platform builds
set -e

echo "🚀 Building Postgirl for all platforms..."

# Clean previous builds
rm -rf dist
mkdir -p dist

# Build flags
LDFLAGS="-s -w"

echo "📦 Building macOS (ARM64)..."
GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -trimpath -ldflags "$LDFLAGS" -o dist/postgirl-darwin-arm64 ./cmd/postgirl

echo "📦 Building macOS (AMD64)..."
GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -trimpath -ldflags "$LDFLAGS" -o dist/postgirl-darwin-amd64 ./cmd/postgirl

echo "📦 Building Linux (AMD64)..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "$LDFLAGS" -o dist/postgirl-linux-amd64 ./cmd/postgirl

echo "📦 Building Linux (ARM64)..."
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags "$LDFLAGS" -o dist/postgirl-linux-arm64 ./cmd/postgirl

echo "📦 Building Windows (AMD64)..."
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "$LDFLAGS" -o dist/postgirl-windows-amd64.exe ./cmd/postgirl

echo "📦 Building Windows (ARM64)..."
GOOS=windows GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags "$LDFLAGS" -o dist/postgirl-windows-arm64.exe ./cmd/postgirl

echo "📁 Copying web interface files..."
cp -r web dist/

echo "✅ Build complete! Single executables ready:"
ls -la dist/

echo ""
echo "🎯 Usage - Just double-click the executable:"
echo ""
echo "📱 macOS:"
echo "  • postgirl-darwin-arm64 (Auto-starts Web UI + opens browser)"
echo "  • postgirl-darwin-arm64 tui (Terminal UI)"
echo ""
echo "🐧 Linux:"
echo "  • ./postgirl-linux-amd64 (Auto-starts Web UI)"
echo "  • ./postgirl-linux-amd64 tui (Terminal UI)"
echo ""
echo "🪟 Windows:"
echo "  • postgirl-windows-amd64.exe (Auto-starts Web UI + opens browser)"
echo "  • postgirl-windows-amd64.exe tui (Terminal UI)"
echo ""
echo "🌐 Web interface will be available at http://localhost:8080"
