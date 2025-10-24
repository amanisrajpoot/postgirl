#!/bin/bash

# Final build script for Postgirl - Completely standalone executables
set -e

echo "🚀 Building Postgirl - Final Standalone Distributions..."

# Clean previous builds
rm -rf dist-final
mkdir -p dist-final

# Build flags
LDFLAGS="-s -w -X main.version=$(git describe --tags --always --dirty)"

# Function to build completely standalone executable
build_standalone() {
    local os=$1
    local arch=$2
    local ext=""
    local cgo_enabled="0"
    
    if [ "$os" = "windows" ]; then
        ext=".exe"
    fi
    
    if [ "$os" = "darwin" ]; then
        cgo_enabled="1"
    fi
    
    local output_name="postgirl-${os}-${arch}${ext}"
    
    echo "📦 Building $os/$arch..."
    
    # Build single executable with interactive menu
    GOOS=$os GOARCH=$arch CGO_ENABLED=$cgo_enabled go build -trimpath -ldflags "$LDFLAGS" -o "dist-final/${output_name}" ./cmd/postgirl
    
    # Create a simple README for each platform
    cat > "dist-final/${os}-${arch}-README.md" << EOF
# Postgirl - $os/$arch

## Quick Start

### Interactive Menu (Default)
\`\`\`bash
./${output_name}
\`\`\`
- Interactive menu to choose interface
- Auto-opens browser for web interface
- No external dependencies

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
- 📁 Collections & Environments
- 📝 Request/Response History
- 🔧 JavaScript Scripting
- 🚀 Cross-platform Support

## Support
- Web Interface: http://localhost:8080
- Terminal Interface: Keyboard navigation
- No installation required - just run!

## Files
- \`${output_name}\` - Postgirl executable (with interactive menu)
- \`${os}-${arch}-README.md\` - This file

EOF
}

# Build for all platforms
build_standalone "darwin" "arm64"
build_standalone "darwin" "amd64"
build_standalone "linux" "amd64"
build_standalone "linux" "arm64"
build_standalone "windows" "amd64"
build_standalone "windows" "arm64"

echo "✅ Final standalone builds complete!"
echo ""
echo "📁 Distribution files created:"
ls -la dist-final/

echo ""
echo "🎯 Usage:"
echo ""
echo "📱 macOS:"
echo "  • postgirl-darwin-arm64 (Interactive menu + direct launch)"
echo ""
echo "🐧 Linux:"
echo "  • ./postgirl-linux-amd64 (Interactive menu + direct launch)"
echo ""
echo "🪟 Windows:"
echo "  • postgirl-windows-amd64.exe (Interactive menu + direct launch)"
echo ""
echo "🌐 Each executable is completely standalone!"
echo "📦 No external files or dependencies required!"
