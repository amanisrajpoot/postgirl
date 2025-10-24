package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/bubbletea"
	"postgirl/internal/app"
	"postgirl/internal/storage"
	"postgirl/internal/storage/sqlite"
	"postgirl/internal/ui"
	"postgirl/internal/web"
)

//go:embed web/static
var webAssets embed.FS

var version = "dev"

func main() {
	var showVersion bool
	var tui bool
	var web bool
	var port int
	
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.BoolVar(&tui, "tui", false, "Start the terminal user interface")
	flag.BoolVar(&web, "web", false, "Start the web interface")
	flag.IntVar(&port, "port", 8080, "Port for web interface")
	
	// Parse flags
	flag.Parse()
	
	// Check for positional arguments
	args := flag.Args()
	if len(args) > 0 {
		switch args[0] {
		case "tui":
			tui = true
		case "web":
			web = true
		}
	}

	if showVersion {
		fmt.Printf("litepost version %s\n", version)
		os.Exit(0)
	}

	if tui {
		// Start the TUI application
		app := ui.NewApp()
		p := tea.NewProgram(app, tea.WithAltScreen())
		
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error running application: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if web {
		// Start the web application
		startWebInterface(port)
		return
	}

	// Default behavior - show interactive menu
	showInteractiveMenu()
}

// startWebInterface starts the web interface
func startWebInterface(port int) {
	fmt.Printf("Starting Postgirl web interface on port %d...\n", port)
	fmt.Printf("Open your browser and go to: http://localhost:%d\n", port)
	fmt.Println("Press Ctrl+C to stop the server")
	
	// Initialize storage
	var storageInstance storage.Storage
	sqliteStorage, err := sqlite.NewSQLiteStorage("postgirl.db")
	if err != nil {
		log.Printf("Warning: Failed to initialize SQLite database: %v", err)
		log.Printf("Falling back to in-memory storage")
		storageInstance = storage.NewMemoryStorage()
	} else {
		storageInstance = sqliteStorage
	}
	
	// Initialize service
	service := app.NewService(storageInstance)
	
	// Create and start web server
	server := web.NewServer(service, port, webAssets)
	
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start web server: %v", err)
	}
}

// showInteractiveMenu shows an interactive menu to choose interface
func showInteractiveMenu() {
	clearScreen()
	
	fmt.Println("🚀 Postgirl")
	fmt.Println("═══════════════════════════════════════")
	fmt.Println("Lightweight Postman Alternative")
	fmt.Println("")
	
	for {
		showMenu()
		
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter your choice (1-3): ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)
		
		switch choice {
		case "1":
			launchWebInterface()
		case "2":
			launchTUI()
		case "3":
			fmt.Println("👋 Goodbye!")
			os.Exit(0)
		default:
			fmt.Println("❌ Invalid choice. Please enter 1, 2, or 3.")
			time.Sleep(1 * time.Second)
		}
	}
}

func showMenu() {
	clearScreen()
	fmt.Println("🚀 Postgirl")
	fmt.Println("═══════════════════════════════════════")
	fmt.Println("")
	fmt.Println("Choose an interface:")
	fmt.Println("")
	fmt.Println("1. 🌐 Web Interface")
	fmt.Println("   • Modern browser-based UI")
	fmt.Println("   • Auto-opens browser")
	fmt.Println("   • Full-featured interface")
	fmt.Println("")
	fmt.Println("2. 💻 Terminal Interface")
	fmt.Println("   • Keyboard-driven TUI")
	fmt.Println("   • No browser required")
	fmt.Println("   • Lightweight and fast")
	fmt.Println("")
	fmt.Println("3. ❌ Exit")
	fmt.Println("")
}

func launchWebInterface() {
	clearScreen()
	fmt.Println("🚀 Starting Postgirl Web Interface...")
	fmt.Println("")
	
	// Start the web server
	cmd := exec.Command(os.Args[0], "web", "--port", "8080")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Start(); err != nil {
		fmt.Printf("❌ Error: Failed to start web interface: %v\n", err)
		waitForKey()
		return
	}

	// Wait a moment for server to start
	fmt.Println("⏳ Starting web server...")
	time.Sleep(2 * time.Second)
	
	// Open browser
	fmt.Println("🌐 Opening browser at http://localhost:8080")
	openBrowser("http://localhost:8080")
	
	fmt.Println("")
	fmt.Println("✅ Web Interface started successfully!")
	fmt.Println("🌐 Browser opened at: http://localhost:8080")
	fmt.Println("📝 Make HTTP requests in your browser")
	fmt.Println("")
	fmt.Println("Press Ctrl+C to stop the server and return to launcher")
	fmt.Println("")
	
	// Keep the process running
	cmd.Wait()
}

func launchTUI() {
	clearScreen()
	fmt.Println("🚀 Starting Postgirl Terminal Interface...")
	fmt.Println("")
	
	// Start the TUI
	cmd := exec.Command(os.Args[0], "tui")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	
	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ Error: Failed to start terminal interface: %v\n", err)
		waitForKey()
		return
	}
	
	fmt.Println("✅ Terminal Interface completed")
	waitForKey()
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default: // linux
		cmd = exec.Command("xdg-open", url)
	}
	
	cmd.Start()
}

func clearScreen() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func waitForKey() {
	fmt.Print("Press Enter to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
