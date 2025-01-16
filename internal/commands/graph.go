package commands

import (
	"fmt"
	"net/http"

	"github.com/baudevs/yolo-cli/internal/web"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func GraphCmd() *cobra.Command {
	var port int

	cmd := &cobra.Command{
		Use:   "graph",
		Short: "Start the interactive graph visualization server",
		Long: `Start a web server that provides an interactive visualization of your project's
epics, features, tasks, and their relationships.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			r := chi.NewRouter()
			r.Use(middleware.Logger)
			r.Use(middleware.Recoverer)

			// Serve static files from embedded filesystem
			fileServer := http.FileServer(http.FS(web.WebFiles))
			r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

			// API routes
			r.Get("/api/nodes", handleGetNodes)
			r.Get("/api/ws", handleWebSocket)

			// Serve index.html for all other routes
			r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
				indexHTML, err := web.WebFiles.ReadFile("static/index.html")
				if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "text/html")
				w.Write(indexHTML)
			})

			addr := fmt.Sprintf(":%d", port)
			fmt.Printf("Starting graph server at http://localhost%s\n", addr)
			return http.ListenAndServe(addr, r)
		},
	}

	cmd.Flags().IntVarP(&port, "port", "p", 4010, "Port to run the server on")
	return cmd
}

func handleGetNodes(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement node data retrieval
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{
		"nodes": [
			{"id": "1", "type": "epic", "title": "Example Epic", "status": "in_progress"},
			{"id": "2", "type": "feature", "title": "Example Feature", "status": "todo"},
			{"id": "3", "type": "task", "title": "Example Task", "status": "done"}
		],
		"links": [
			{"source": "1", "target": "2"},
			{"source": "2", "target": "3"}
		]
	}`))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Error upgrading to websocket: %v\n", err)
		return
	}
	defer conn.Close()

	// TODO: Implement real-time updates
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		conn.WriteMessage(messageType, message)
	}
} 