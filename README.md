# 🚀 Litepost

A lightweight Postman alternative built with Go, featuring both Terminal UI and Web UI.

## ✨ Features

- **Dual Interface**: Terminal UI (TUI) and Web UI
- **HTTP Methods**: GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS
- **Authentication**: Basic Auth, Bearer Token, API Key, OAuth2, Digest, Hawk
- **Request History**: Automatic saving of requests and responses
- **Collections**: Organize requests into collections
- **Environments**: Variable management across requests
- **Scripting**: Pre-request and post-response JavaScript scripts
- **Cross-platform**: macOS, Linux, Windows (AMD64 & ARM64)

## 🚀 Quick Start

### Build All Platforms
```bash
./build.sh
```

### Run Web Interface
```bash
# macOS
./dist/litepost-darwin-arm64 web

# Linux
./dist/litepost-linux-amd64 web

# Windows
dist\litepost-windows-amd64.exe web
```

### Run Terminal Interface
```bash
# macOS
./dist/litepost-darwin-arm64 tui

# Linux
./dist/litepost-linux-amd64 tui

# Windows
dist\litepost-windows-amd64.exe tui
```

## 🌐 Web Interface

Open your browser to `http://localhost:8080` after starting the web interface.

## 📦 Available Executables

- `litepost-darwin-arm64` - macOS (Apple Silicon)
- `litepost-darwin-amd64` - macOS (Intel)
- `litepost-linux-amd64` - Linux (AMD64)
- `litepost-linux-arm64` - Linux (ARM64)
- `litepost-windows-amd64.exe` - Windows (AMD64)
- `litepost-windows-arm64.exe` - Windows (ARM64)

## 🔧 Development

```bash
# Install dependencies
go mod download

# Run in development
go run ./cmd/litepost web
```

## 📝 License

MIT License
