/*
================================================================================
MAIN ENTRY POINT - P2P Academic Library "The Knowledge Exchange"
================================================================================
This is the main entry point for the Knowledge Exchange P2P Academic Library.

Go Concepts Used:
- Main function: Application entry point
- Packages: Module organization
- Goroutines: Concurrent service startup
- Signals: Graceful shutdown handling
================================================================================
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"knowledge-exchange/gateway"
	"knowledge-exchange/utils"
)

// ============================================================================
// CONSTANTS
// ============================================================================

const (
	banner = `
╔═══════════════════════════════════════════════════════════════╗
║                                                               ║
║        📚 THE KNOWLEDGE EXCHANGE 📚                           ║
║        P2P Academic Library                                   ║
║                                                               ║
║        Decentralized • Fair • Open                            ║
║                                                               ║
╚═══════════════════════════════════════════════════════════════╝
`
)

// ============================================================================
// MAIN FUNCTION
// ============================================================================

func main() {
	// Print banner
	fmt.Println(banner)
	fmt.Printf("Version: %s\n\n", utils.AppVersion)

	// Parse command line flags
	var (
		port    = flag.Int("port", utils.DefaultAPIPort, "API server port")
		name    = flag.String("name", "Anonymous Peer", "Peer display name")
		dataDir = flag.String("data", utils.DefaultDataDir, "Data storage directory")
	)
	flag.Parse()

	// Log startup
	log.Printf("Starting Knowledge Exchange with configuration:")
	log.Printf("  - Port: %d", *port)
	log.Printf("  - Name: %s", *name)
	log.Printf("  - Data Directory: %s", *dataDir)

	// Create configuration
	config := utils.DefaultConfig()
	config.APIPort = *port
	config.PeerName = *name
	config.DataDir = *dataDir
	config.SharedFilesDir = *dataDir + "/sharedFiles"
	config.TempDir = *dataDir + "/temp"

	// Ensure directories exist
	if err := utils.EnsureDirectories(); err != nil {
		log.Fatalf("Failed to create directories: %v", err)
	}
	log.Println("✓ Directories initialized")

	// Generate peer ID
	localIP, _ := utils.GetLocalIP()
	if localIP == "" {
		localIP = "127.0.0.1"
	}
	config.PeerID = utils.GeneratePeerID(config.PeerName, localIP, *port)
	config.HostIP = localIP
	log.Printf("✓ Peer ID: %s", config.PeerID)

	// Create and start server
	server := gateway.NewServer(config)

	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	log.Printf("✓ Server started on http://%s:%d", localIP, *port)

	// Print helpful information
	fmt.Println("\n" + strings.Repeat("─", 60))
	fmt.Println("API Endpoints:")
	fmt.Println("  GET  /api/health         - Health check")
	fmt.Println("  GET  /api/status         - Server status")
	fmt.Println("  POST /api/peers/register - Register as peer")
	fmt.Println("  GET  /api/peers          - List all peers")
	fmt.Println("  GET  /api/files          - List all files")
	fmt.Println("  GET  /api/files/search   - Search files (?q=query)")
	fmt.Println("  POST /api/files/upload   - Upload file")
	fmt.Println("  GET  /api/files/download - Download file")
	fmt.Println("  GET  /api/reputation     - Get reputation")
	fmt.Println("  POST /api/ratings/file   - Rate a file")
	fmt.Println("  GET  /api/stats          - System statistics")
	fmt.Println(strings.Repeat("─", 60))
	fmt.Println("\nPress Ctrl+C to stop the server\n")

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	log.Println("\nShutting down server...")
	if err := server.Stop(); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}

	log.Println("Server stopped. Goodbye! 👋")
}
