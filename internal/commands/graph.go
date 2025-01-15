package commands

import (
	"fmt"
	"path/filepath"

	"github.com/baudevs/yolo-cli/internal/server"
	"github.com/spf13/cobra"
)

func GraphCmd() *cobra.Command {
	var port int

	cmd := &cobra.Command{
		Use:   "graph",
		Short: "Start the interactive project graph visualization",
		Long: `Start a web server that provides an interactive visualization of your project's
epics, features, tasks, and their relationships. The graph is rendered using WebGL
and provides real-time updates through WebSocket connections.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get the web directory path relative to the current directory
			webDir := filepath.Join("web")

			// Create and start the graph server
			srv := server.NewGraphServer(webDir, port)
			
			// Load project data
			if err := srv.LoadProjectData("."); err != nil {
				return fmt.Errorf("failed to load project data: %w", err)
			}

			// Start the server
			return srv.Start()
		},
	}

	// Add flags
	cmd.Flags().IntVarP(&port, "port", "p", 4010, "Port to run the graph server on")

	return cmd
} 