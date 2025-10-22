#!/usr/bin/env bash
set -euo pipefail
REPO="yourname/litepost"
VERSION="${1:-latest}"

detect() {
  OS=$(uname -s | tr '[:upper:]' '[:lower:]')
  ARCH=$(uname -m)
  case "$ARCH" in
    x86_64|amd64) ARCH=amd64 ;;
    aarch64|arm64) ARCH=arm64 ;;
    *) echo "Unsupported arch: $ARCH"; exit 1 ;;
  esac
  echo "$OS" "$ARCH"
}

dl() {
  local os="$1" arch="$2" tmp="$(mktemp -d)"
  cd "$tmp"
  if [ "$os" = "windows" ]; then
    echo "Use PowerShell installer on Windows."
    exit 1
  fi
  if [ "$VERSION" = "latest" ]; then
    TAG=$(curl -s https://api.github.com/repos/$REPO/releases/latest | grep -Po '"tag_name": "\K.*?(?=")')
  else
    TAG="$VERSION"
  fi
  FILE="litepost-${os}-${arch}.$([ "$os" = "linux" ] && echo tar.gz || echo zip)"
  URL="https://github.com/$REPO/releases/download/$TAG/$FILE"
  echo "Downloading $URL"
  curl -L -o "$FILE" "$URL"
  mkdir -p "$HOME/.local/bin"
  if [[ "$FILE" == *.tar.gz ]]; then tar -xzf "$FILE"; else unzip -q "$FILE"; fi
  cp -f litepost-${os}-${arch}/litepost "$HOME/.local/bin/litepost"
  chmod +x "$HOME/.local/bin/litepost"
  echo "Installed to $HOME/.local/bin. Ensure it's on PATH."
}

read OS ARCH < <(detect)
dl "$OS" "$ARCH"
