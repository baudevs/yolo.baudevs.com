package shortcuts

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

// Client represents a WebSocket client
type Client struct {
	conn *websocket.Conn
	send chan []byte
}

// Server represents the shortcuts server
type Server struct {
	daemon     *Daemon
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

// NewServer creates a new shortcuts server
func NewServer(daemon *Daemon) *Server {
	s := &Server{
		daemon:     daemon,
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

	// Set up shortcut trigger handler
	daemon.OnShortcut(func(shortcut *Shortcut) {
		// Show native notification
		go func() {
			title := fmt.Sprintf("YOLO Shortcut: %s", shortcut.Name)
			message := fmt.Sprintf("Executing command: %s", shortcut.Command)
			notification := fmt.Sprintf(`display notification "%s" with title "%s" sound name "Glass"`, message, title)
			
			cmd := exec.Command("osascript", "-e", notification)
			if err := cmd.Run(); err != nil {
				log.Printf("Failed to show notification: %v", err)
			}
		}()

		// Execute the command
		go func() {
			cmd := exec.Command("yolo", shortcut.Command)
			if err := cmd.Start(); err != nil {
				log.Printf("Failed to execute command: %v", err)
			}
		}()

		// Notify clients
		msg := map[string]interface{}{
			"type":     "shortcut_triggered",
			"shortcut": shortcut,
		}
		if data, err := json.Marshal(msg); err == nil {
			s.broadcast <- data
		}
	})

	go s.run()
	return s
}

func (s *Server) run() {
	for {
		select {
		case client := <-s.register:
			s.clients[client] = true
			log.Printf("Client connected. Total clients: %d", len(s.clients))

		case client := <-s.unregister:
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				close(client.send)
				log.Printf("Client disconnected. Total clients: %d", len(s.clients))
			}

		case message := <-s.broadcast:
			for client := range s.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(s.clients, client)
				}
			}
		}
	}
}

// Routes returns the router for shortcuts endpoints
func (s *Server) Routes() chi.Router {
	r := chi.NewRouter()

	// Static files
	r.Handle("/*", http.FileServer(http.Dir("internal/web/static/shortcuts")))

	// API endpoints
	r.Route("/api", func(r chi.Router) {
		r.Get("/list", s.handleList)
		r.Post("/save", s.handleSave)
		r.Post("/delete", s.handleDelete)
		r.Post("/toggle", s.handleToggle)
	})

	// WebSocket endpoint
	r.Get("/ws", s.handleWebSocket)

	return r
}

func (s *Server) handleList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	shortcuts := s.daemon.GetShortcuts()
	if err := json.NewEncoder(w).Encode(shortcuts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleSave(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var shortcut Shortcut
	if err := json.NewDecoder(r.Body).Decode(&shortcut); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.daemon.AddShortcut(&shortcut); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Notify clients
	msg := map[string]string{"type": "shortcuts_updated"}
	if data, err := json.Marshal(msg); err == nil {
		s.broadcast <- data
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.daemon.RemoveShortcut(req.ID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Notify clients
	msg := map[string]string{"type": "shortcuts_updated"}
	if data, err := json.Marshal(msg); err == nil {
		s.broadcast <- data
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (s *Server) handleToggle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req struct {
		ID      string `json:"id"`
		Enabled bool   `json:"enabled"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.daemon.SetEnabled(req.ID, req.Enabled); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Notify clients
	msg := map[string]string{"type": "shortcuts_updated"}
	if data, err := json.Marshal(msg); err == nil {
		s.broadcast <- data
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	// Register client
	client := &Client{conn: conn, send: make(chan []byte, 256)}
	s.register <- client
	defer func() {
		s.unregister <- client
	}()

	// Send initial shortcuts list
	shortcuts := s.daemon.GetShortcuts()
	if data, err := json.Marshal(map[string]interface{}{
		"type": "shortcuts_list",
		"shortcuts": shortcuts,
	}); err == nil {
		client.send <- data
	}

	// Start client goroutines
	go client.writePump(s)
	client.readPump(s)
}

func (c *Client) writePump(s *Server) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) readPump(s *Server) {
	defer func() {
		s.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			break
		}

		// Handle incoming messages if needed
		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err == nil {
			// Process message based on type
			switch msg["type"] {
			case "request_shortcuts":
				shortcuts := s.daemon.GetShortcuts()
				if data, err := json.Marshal(map[string]interface{}{
					"type": "shortcuts_list",
					"shortcuts": shortcuts,
				}); err == nil {
					c.send <- data
				}
			}
		}
	}
} 