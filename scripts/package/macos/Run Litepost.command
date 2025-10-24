#!/usr/bin/env bash
cd "$(dirname "$0")"
chmod +x ./postgirl
# Open a new Terminal window and run TUI (or change to 'send'/'run')
osascript <<'APPLESCRIPT'
tell application "Terminal"
  activate
  do script "cd \"$(PWD)\" && ./postgirl tui || read -n 1 -s -r -p \"Press any key to close...\""
end tell
APPLESCRIPT
