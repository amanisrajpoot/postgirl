package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/charmbracelet/bubbletea"
	"postgirl/internal/app"
	"postgirl/internal/storage/sqlite"
	"postgirl/internal/ui"
	"postgirl/internal/web"
)

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

	// Default behavior - start web interface automatically
	fmt.Println("🚀 Starting Postgirl Web Interface...")
	fmt.Println("Opening browser at http://localhost:8080")
	
	// Open browser automatically on macOS
	go func() {
		time.Sleep(2 * time.Second)
		exec.Command("open", "http://localhost:8080").Run()
	}()
	
	startWebInterface(port)
}

// startWebInterface starts the web interface
func startWebInterface(port int) {
	fmt.Printf("Starting Postgirl web interface on port %d...\n", port)
	fmt.Printf("Open your browser and go to: http://localhost:%d\n", port)
	fmt.Println("Press Ctrl+C to stop the server")
	
	// Initialize storage
	storage, err := sqlite.NewSQLiteStorage("postgirl.db")
	if err != nil {
		log.Printf("Warning: Failed to initialize database: %v", err)
		storage = nil
	}
	
	// Initialize service
	service := app.NewService(storage)
	
	// Create and start web server
	server := web.NewServer(service, port)
	
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start web server: %v", err)
	}
}
