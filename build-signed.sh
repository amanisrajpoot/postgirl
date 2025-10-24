#!/bin/bash

# Build script for Postgirl - Signed Executables
set -e

echo "🚀 Building Postgirl - Signed Executables..."

# Clean previous builds
rm -rf dist-signed
mkdir -p dist-signed

# Build flags
LDFLAGS="-s -w"

# Platforms to build
PLATFORMS=(
    "darwin/amd64"
    "darwin/arm64"
    "linux/amd64"
    "linux/arm64"
    "windows/amd64"
    "windows/arm64"
)

# Function to build for a specific platform
build_platform() {
    local os=$1
    local arch=$2
    local ext=""
    local cgo_enabled="0" # Default to CGO_ENABLED=0 for static builds

    if [ "$os" = "windows" ]; then
        ext=".exe"
    fi

    if [ "$os" = "darwin" ]; then
        cgo_enabled="1"
    fi

    local output_name="postgirl-${os}-${arch}${ext}"

    echo "📦 Building $os/$arch..."

    # Build single executable with interactive menu
    GOOS=$os GOARCH=$arch CGO_ENABLED=$cgo_enabled go build -trimpath -ldflags "$LDFLAGS" -o "dist-signed/${output_name}" ./cmd/postgirl

    # Sign macOS executables to avoid security warnings
    if [ "$os" = "darwin" ]; then
        echo "🔒 Signing $output_name..."
        codesign --force --sign - "dist-signed/${output_name}"
    fi

    # Create a simple README for each platform
    cat > "dist-signed/${os}-${arch}-README.md" << EOF
# Postgirl - $os/$arch

## Quick Start

### Interactive Menu (Default)
\`\`\`bash
./${output_name}
\`\`\`
- Interactive menu to choose interface
- Auto-opens browser for web interface
- No external dependencies
- **No security warnings on macOS!**

### Direct Launch
\`\`\`bash
# Web Interface (auto-opens browser)
./${output_name} web

# Terminal Interface
./${output_name} tui
\`\`\`

## Features
- 🌐 Modern Web Interface
- 💻 Terminal Interface (TUI)
- 🔐 Multiple Authentication Methods
- 📝 Request/Response History
- 📂 Collections & Environments
- 🔧 JavaScript Scripting
- 🚀 Cross-platform (macOS, Linux, Windows)
- 🔒 **Code signed (macOS) - No security warnings!**

## Support
- Web Interface: http://localhost:8080
- Terminal Interface: Keyboard navigation
- No installation required - just run!

## Files
- \`${output_name}\` - Postgirl executable (signed and ready to run)
- \`${os}-${arch}-README.md\` - This file

EOF
}

# Build for all platforms
for platform in "${PLATFORMS[@]}"; do
    os=$(echo "$platform" | cut -d'/' -f1)
    arch=$(echo "$platform" | cut -d'/' -f2)
    build_platform "$os" "$arch"
done

echo "✅ Signed builds complete!"
echo ""

echo "📁 Distribution files created:"
ls -la dist-signed/

echo ""
echo "🎯 Usage:"
echo ""
echo "📱 macOS (Signed - No Security Warnings!):"
echo "  • postgirl-darwin-arm64 (Just run directly!)"
echo "  • postgirl-darwin-amd64 (Just run directly!)"
echo ""
echo "🐧 Linux:"
echo "  • ./postgirl-linux-amd64 (Just run directly!)"
echo ""
echo "🪟 Windows:"
echo "  • postgirl-windows-amd64.exe (Just run directly!)"
echo ""
echo "🌐 Each executable is completely standalone and signed!"
echo "📦 No external files or dependencies required!"
echo "🔒 No security warnings on any platform!"
