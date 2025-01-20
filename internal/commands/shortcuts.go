package commands

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/baudevs/yolo.baudevs.com/internal/shortcuts"
	"github.com/spf13/cobra"
)

var ShortcutsCmd = &cobra.Command{
	Use:   "shortcuts",
	Short: "⌨️  Configure system-wide shortcuts",
	Long: `🎮 Configure and manage system-wide keyboard shortcuts.

This command starts a web interface where you can:
1. Configure global shortcuts
2. Test keyboard combinations
3. Enable/disable shortcuts
4. View active shortcuts

Note: This feature is currently in development.
Some functionality may be limited.`,
	RunE: runShortcuts,
}

func runShortcuts(cmd *cobra.Command, args []string) error {
	fmt.Println("🚀 Starting shortcuts configuration...")

	// Create server
	server, err := shortcuts.NewServer()
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	// Start server
	if err := server.Start(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	// Set up routes
	http.HandleFunc("/ws", server.HandleWebSocket)
	http.Handle("/", http.FileServer(http.Dir("internal/web/static/shortcuts")))

	// Start HTTP server
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Printf("Failed to start server: %v\n", err)
			os.Exit(1)
		}
	}()

	fmt.Println("🌐 Web interface available at http://localhost:8080")
	fmt.Println("🛑 Press Ctrl+C to stop")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\n👋 Shutting down...")
	return server.Stop()
} 