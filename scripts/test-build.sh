#!/usr/bin/env bash
# Test script to verify the build system works

set -e

echo "Testing Litepost build system..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go 1.22 or later."
    echo "Visit: https://golang.org/dl/"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | cut -d' ' -f3 | sed 's/go//')
echo "Go version: $GO_VERSION"

# Initialize Go module
echo "Initializing Go module..."
go mod tidy

# Test basic build
echo "Testing basic build..."
go build -o /tmp/litepost-test ./cmd/litepost

# Test the binary
echo "Testing binary..."
/tmp/litepost-test --version
/tmp/litepost-test --help

# Clean up
rm -f /tmp/litepost-test

echo "✅ Basic build test passed!"

# Test Makefile (if make is available)
if command -v make &> /dev/null; then
    echo "Testing Makefile..."
    make clean
    echo "✅ Makefile test passed!"
else
    echo "⚠️  Make not found - skipping Makefile test"
fi

echo "🎉 All tests passed! Your build system is ready."
