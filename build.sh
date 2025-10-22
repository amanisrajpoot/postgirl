#!/bin/bash

# Build script for Litepost - Cross-platform builds
set -e

echo "🚀 Building Litepost for all platforms..."

# Clean previous builds
rm -rf dist
mkdir -p dist

# Build flags
LDFLAGS="-s -w"

echo "📦 Building macOS (ARM64)..."
GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -trimpath -ldflags "$LDFLAGS" -o dist/litepost-darwin-arm64 ./cmd/litepost

echo "📦 Building macOS (AMD64)..."
GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -trimpath -ldflags "$LDFLAGS" -o dist/litepost-darwin-amd64 ./cmd/litepost

echo "📦 Building Linux (AMD64)..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "$LDFLAGS" -o dist/litepost-linux-amd64 ./cmd/litepost

echo "📦 Building Linux (ARM64)..."
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags "$LDFLAGS" -o dist/litepost-linux-arm64 ./cmd/litepost

echo "📦 Building Windows (AMD64)..."
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "$LDFLAGS" -o dist/litepost-windows-amd64.exe ./cmd/litepost

echo "📦 Building Windows (ARM64)..."
GOOS=windows GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags "$LDFLAGS" -o dist/litepost-windows-arm64.exe ./cmd/litepost

echo "📁 Copying web interface files..."
cp -r web dist/

echo "✅ Build complete! Single executables ready:"
ls -la dist/

echo ""
echo "🎯 Usage - Just double-click the executable:"
echo ""
echo "📱 macOS:"
echo "  • litepost-darwin-arm64 (Auto-starts Web UI + opens browser)"
echo "  • litepost-darwin-arm64 tui (Terminal UI)"
echo ""
echo "🐧 Linux:"
echo "  • ./litepost-linux-amd64 (Auto-starts Web UI)"
echo "  • ./litepost-linux-amd64 tui (Terminal UI)"
echo ""
echo "🪟 Windows:"
echo "  • litepost-windows-amd64.exe (Auto-starts Web UI + opens browser)"
echo "  • litepost-windows-amd64.exe tui (Terminal UI)"
echo ""
echo "🌐 Web interface will be available at http://localhost:8080"
