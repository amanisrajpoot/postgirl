@echo off
cd /d "%~dp0"
set BIN=litepost.exe
if not exist "%BIN%" (
  echo %BIN% not found.
  pause
  exit /b 1
)
"%BIN%" tui
echo.
echo Press any key to close...
pause >nul
