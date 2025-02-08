package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
)

type GraphServer struct {
	router     *chi.Mux
	wsHub      *WSHub
	dataLoader *ProjectDataLoader
	staticPath string
	port       int
}

type ProjectDataLoader struct {
	sync.RWMutex
	nodes map[string]*GraphNode
	links map[string][]string
}

type GraphNode struct {
	ID       string    `json:"id"`
	Type     string    `json:"type"` // epic, feature, task
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	Links    []string  `json:"links"`
	Status   string    `json:"status"`
	Version  string    `json:"version"`
	Modified time.Time `json:"modified"`
}

// WebSocket hub for real-time updates
type WSHub struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // In production, implement proper origin checking
	},
}

func NewGraphServer(staticPath string, port int) *GraphServer {
	hub := &WSHub{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}

	go hub.run()

	return &GraphServer{
		router:     chi.NewRouter(),
		wsHub:      hub,
		dataLoader: &ProjectDataLoader{
			nodes: make(map[string]*GraphNode),
			links: make(map[string][]string),
		},
		staticPath: staticPath,
		port:       port,
	}
}

func (s *GraphServer) setupRoutes() {
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Compress(5))

	// API routes
	s.router.Route("/api", func(r chi.Router) {
		r.Get("/nodes", s.handleGetNodes)
		r.Get("/node/{id}", s.handleGetNode)
		r.Get("/links", s.handleGetLinks)
		r.Get("/ws", s.handleWebSocket)
	})

	// Static file serving with proper MIME types
	fileServer := http.FileServer(http.Dir(s.staticPath))
	s.router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		// Set correct MIME types
		ext := filepath.Ext(r.URL.Path)
		switch ext {
		case ".css":
			w.Header().Set("Content-Type", "text/css")
		case ".js":
			w.Header().Set("Content-Type", "application/javascript")
		case ".html":
			w.Header().Set("Content-Type", "text/html")
		}
		fileServer.ServeHTTP(w, r)
	})
}

func (s *GraphServer) Start() error {
	s.setupRoutes()
	addr := fmt.Sprintf(":%d", s.port)
	fmt.Printf("Starting graph server on http://localhost%s\n", addr)
	return http.ListenAndServe(addr, s.router)
}

func (s *GraphServer) handleGetNodes(w http.ResponseWriter, r *http.Request) {
	s.dataLoader.RLock()
	defer s.dataLoader.RUnlock()

	nodes := make([]*GraphNode, 0, len(s.dataLoader.nodes))
	for _, node := range s.dataLoader.nodes {
		nodes = append(nodes, node)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(nodes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *GraphServer) handleGetNode(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	s.dataLoader.RLock()
	node, exists := s.dataLoader.nodes[id]
	s.dataLoader.RUnlock()

	if !exists {
		http.Error(w, "Node not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(node)
}

func (s *GraphServer) handleGetLinks(w http.ResponseWriter, r *http.Request) {
	s.dataLoader.RLock()
	defer s.dataLoader.RUnlock()

	links := make([]map[string]string, 0)
	for nodeID, nodeLinks := range s.dataLoader.links {
		for _, targetID := range nodeLinks {
			links = append(links, map[string]string{
				"source": nodeID,
				"target": targetID,
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(links); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *GraphServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("WebSocket upgrade error: %v\n", err)
		return
	}

	s.wsHub.register <- conn

	// Handle incoming messages
	go func() {
		defer func() {
			s.wsHub.unregister <- conn
			conn.Close()
		}()

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}
			// Handle message (e.g., search requests, filters)
			s.handleWSMessage(message)
		}
	}()
}

func (s *GraphServer) handleWSMessage(message []byte) {
	// TODO: Implement message handling
	// Example: search requests, filter updates, etc.
}

func (h *WSHub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.Close()
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				err := client.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					client.Close()
					delete(h.clients, client)
				}
			}
		}
	}
}

func (s *GraphServer) LoadProjectData(projectPath string) error {
	s.dataLoader.Lock()
	defer s.dataLoader.Unlock()

	// Load epics
	epicsPath := filepath.Join(projectPath, "yolo/epics")
	if err := s.loadNodesFromDir(epicsPath, "epic"); err != nil {
		return fmt.Errorf("failed to load epics: %w", err)
	}

	// Load features
	featuresPath := filepath.Join(projectPath, "yolo/features")
	if err := s.loadNodesFromDir(featuresPath, "feature"); err != nil {
		return fmt.Errorf("failed to load features: %w", err)
	}

	// Load tasks
	tasksPath := filepath.Join(projectPath, "yolo/tasks")
	if err := s.loadNodesFromDir(tasksPath, "task"); err != nil {
		return fmt.Errorf("failed to load tasks: %w", err)
	}

	return nil
}

func (s *GraphServer) loadNodesFromDir(dirPath string, nodeType string) error {
	// Create some test nodes for now
	nodes := []struct {
		id      string
		title   string
		content string
		links   []string
	}{
		{
			id:      fmt.Sprintf("%s-1", nodeType),
			title:   fmt.Sprintf("%s One", nodeType),
			content: "This is a test node",
			links:   []string{fmt.Sprintf("%s-2", nodeType)},
		},
		{
			id:      fmt.Sprintf("%s-2", nodeType),
			title:   fmt.Sprintf("%s Two", nodeType),
			content: "This is another test node",
			links:   []string{fmt.Sprintf("%s-1", nodeType)},
		},
	}

	for _, n := range nodes {
		node := &GraphNode{
			ID:       n.id,
			Type:     nodeType,
			Title:    n.title,
			Content:  n.content,
			Links:    n.links,
			Status:   "active",
			Version:  "0.1.0",
			Modified: time.Now(),
		}
		s.dataLoader.nodes[node.ID] = node
		s.dataLoader.links[node.ID] = n.links
	}

	return nil
} 