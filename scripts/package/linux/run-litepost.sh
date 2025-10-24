#!/usr/bin/env bash
set -e
cd "$(dirname "$0")"
chmod +x ./postgirl
# Try to open in a terminal emulator if double-clicked from GUI
TERM_EXEC=$(command -v x-terminal-emulator || command -v gnome-terminal || command -v konsole || command -v xfce4-terminal || true)
if [ -n "$TERM_EXEC" ] && [ -z "$PS1" ]; then
  "$TERM_EXEC" -e bash -lc "./postgirl tui; echo; read -p 'Press ENTER to close...'"
else
  ./postgirl tui
fi
