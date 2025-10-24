# 🚀 Postgirl

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
- **Standalone**: Single executable files with no dependencies

## 🚀 Quick Start

### Build All Platforms
```bash
./build-final.sh
```

### Run Web Interface
```bash
# macOS
./dist-final/postgirl-darwin-arm64

# Linux
./dist-final/postgirl-linux-amd64

# Windows
dist-final\postgirl-windows-amd64.exe
```

### Run Terminal Interface
```bash
# macOS
./dist-final/postgirl-darwin-arm64 tui

# Linux
./dist-final/postgirl-linux-amd64 tui

# Windows
dist-final\postgirl-windows-amd64.exe tui
```

### Interactive Launcher (Recommended)
```bash
# macOS
./dist-final/postgirl-launcher-darwin-arm64

# Linux
./dist-final/postgirl-launcher-linux-amd64

# Windows
dist-final\postgirl-launcher-windows-amd64.exe
```

## 🌐 Web Interface

Open your browser to `http://localhost:8080` after starting the web interface.

## 📦 Available Executables

- `postgirl-darwin-arm64` - macOS (Apple Silicon)
- `postgirl-darwin-amd64` - macOS (Intel)
- `postgirl-linux-amd64` - Linux (AMD64)
- `postgirl-linux-arm64` - Linux (ARM64)
- `postgirl-windows-amd64.exe` - Windows (AMD64)
- `postgirl-windows-arm64.exe` - Windows (ARM64)

## 🔧 Development

```bash
# Install dependencies
go mod download

# Run in development
go run ./cmd/postgirl web
```

## 📝 License

MIT License