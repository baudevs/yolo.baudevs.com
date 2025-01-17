package commands

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/baudevs/yolo-cli/internal/shortcuts"
	"github.com/spf13/cobra"
)

// ShortcutsCmd represents the shortcuts command
var ShortcutsCmd = &cobra.Command{
	Use:   "shortcuts",
	Short: "⌨️  Configure system-wide keyboard shortcuts",
	Long: `Configure system-wide keyboard shortcuts for YOLO commands.

This command opens a web interface where you can:
- Add new keyboard shortcuts
- Edit existing shortcuts
- Enable/disable shortcuts
- Delete shortcuts

The shortcuts will work system-wide, allowing you to trigger YOLO commands
from anywhere on your system.

Note: This feature requires accessibility permissions on macOS.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get config directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}

		configDir := filepath.Join(homeDir, ".config", "yolo")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}

		// Create daemon
		daemon, err := shortcuts.NewDaemon(filepath.Join(configDir, "shortcuts.json"))
		if err != nil {
			return fmt.Errorf("failed to create shortcuts daemon: %w", err)
		}

		// Start daemon
		if err := daemon.Start(); err != nil {
			return fmt.Errorf("failed to start shortcuts daemon: %w", err)
		}
		defer daemon.Stop()

		// Create and start server
		server := shortcuts.NewServer(daemon)

		// Get port from flag or use default
		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			port = 4011 // Default port
		}

		fmt.Printf("🎮 Starting shortcuts configuration...\n")
		fmt.Printf("🌟 Configuration interface available at http://localhost:%d\n", port)
		fmt.Printf("⚡️ Quick tips:\n")
		fmt.Printf("🖱️  Click 'Add Shortcut' to create a new shortcut\n")
		fmt.Printf("⌨️  Click the keys field to record a key combination\n")
		fmt.Printf("💾 Changes are saved automatically\n")

		// Start HTTP server
		if err := startServer(port, server.Routes()); err != nil {
			return fmt.Errorf("server error: %w", err)
		}

		return nil
	},
}

func init() {
	ShortcutsCmd.Flags().Int("port", 4011, "Port to run the shortcuts configuration server on")
}

func startServer(port int, handler http.Handler) error {
	addr := fmt.Sprintf(":%d", port)
	return http.ListenAndServe(addr, handler)
} 