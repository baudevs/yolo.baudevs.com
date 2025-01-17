package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/baudevs/yolo-cli/internal/web"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

// Node represents a project item in the graph
type Node struct {
	ID          string   `json:"id"`
	Type        string   `json:"type"`
	Title       string   `json:"title"`
	Status      string   `json:"status"`
	Description string   `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

// Link represents a connection between nodes
type Link struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Type   string `json:"type,omitempty"`
}

// GraphData represents the complete graph structure
type GraphData struct {
	Nodes []Node `json:"nodes"`
	Links []Link `json:"links"`
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// activeConnections keeps track of all WebSocket connections
	activeConnections = struct {
		sync.RWMutex
		conns []*websocket.Conn
	}{}
)

var GraphCmd = &cobra.Command{
	Use:   "graph",
	Short: "🎮 See your project in amazing 3D!",
	Long: `🌈 Experience your project like never before! 

This magical command creates an awesome 3D visualization where you can:
🎯 See all your ideas floating in space
🔍 Click and explore each piece
🎨 Watch how everything connects
🎮 Zoom, spin, and play around
🔍 Search for anything you need

It's like a video game for your project! Perfect for:
👩‍💼 Showing progress to your team
📊 Creating cool presentations
🎓 Understanding how things fit together
🎯 Finding what to work on next

Just run the command and open the link in your browser - no technical skills needed!

The graph will automatically update as you make changes, so you can:
✨ Watch your project grow in real-time
🔄 See new connections form
📈 Track progress visually
🎯 Spot areas that need attention`,
	RunE: runGraph,
}

func runGraph(cmd *cobra.Command, args []string) error {
	fmt.Println("🎮 Preparing your 3D project experience...")

	// Create router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Serve static files from embedded filesystem with correct MIME types
	r.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/static/")
		
		// Set correct MIME types
		if strings.HasSuffix(path, ".js") {
			w.Header().Set("Content-Type", "application/javascript")
		} else if strings.HasSuffix(path, ".css") {
			w.Header().Set("Content-Type", "text/css")
		} else if strings.HasSuffix(path, ".html") {
			w.Header().Set("Content-Type", "text/html")
		}

		content, err := web.WebFiles.ReadFile("static/" + path)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		w.Write(content)
	})

	// API routes
	r.Get("/api/nodes", handleGetNodes)
	r.Get("/ws", handleWebSocket)

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

	// Start server
	port := 4010
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("\n🌟 Your project visualization is ready at http://localhost%s\n", addr)
	fmt.Println("\n�� Quick tips:")
	fmt.Println("🖱️  Left click + drag to rotate")
	fmt.Println("🔍 Scroll to zoom in/out")
	fmt.Println("👆 Click on nodes to see details")
	fmt.Println("🔍 Use the search bar to find anything")

	return http.ListenAndServe(addr, r)
}

func handleGetNodes(w http.ResponseWriter, r *http.Request) {
	// Load project data
	data := loadProjectData()

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Error upgrading to websocket: %v\n", err)
		return
	}

	// Add connection to active connections
	activeConnections.Lock()
	activeConnections.conns = append(activeConnections.conns, conn)
	activeConnections.Unlock()

	defer func() {
		conn.Close()
		// Remove connection from active connections
		activeConnections.Lock()
		for i, c := range activeConnections.conns {
			if c == conn {
				activeConnections.conns = append(activeConnections.conns[:i], activeConnections.conns[i+1:]...)
				break
			}
		}
		activeConnections.Unlock()
	}()

	// Handle incoming messages
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		// Broadcast message to all connections
		activeConnections.RLock()
		for _, c := range activeConnections.conns {
			if err := c.WriteMessage(messageType, message); err != nil {
				fmt.Printf("Error broadcasting message: %v\n", err)
			}
		}
		activeConnections.RUnlock()
	}
}

func loadProjectData() GraphData {
	// TODO: Load actual project data from YOLO files
	// For now, return example data
	return GraphData{
		Nodes: []Node{
			{
				ID:          "E001",
				Type:        "epic",
				Title:       "Project Initialization",
				Status:      "in_progress",
				Description: "Set up the core project structure and features",
				Tags:        []string{"core", "setup"},
			},
			{
				ID:          "F001",
				Type:        "feature",
				Title:       "Interactive CLI",
				Status:      "done",
				Description: "Create an engaging command-line interface",
				Tags:        []string{"ui", "core"},
			},
			{
				ID:          "F002",
				Type:        "feature",
				Title:       "AI Integration",
				Status:      "in_progress",
				Description: "Integrate AI capabilities for enhanced functionality",
				Tags:        []string{"ai", "core"},
			},
			{
				ID:          "T001",
				Type:        "task",
				Title:       "Setup Project Structure",
				Status:      "done",
				Description: "Initialize the basic project structure",
			},
			{
				ID:          "T002",
				Type:        "task",
				Title:       "Implement Graph View",
				Status:      "in_progress",
				Description: "Create the 3D graph visualization",
			},
		},
		Links: []Link{
			{Source: "E001", Target: "F001", Type: "contains"},
			{Source: "E001", Target: "F002", Type: "contains"},
			{Source: "F001", Target: "T001", Type: "implements"},
			{Source: "F002", Target: "T002", Type: "implements"},
		},
	}
} 