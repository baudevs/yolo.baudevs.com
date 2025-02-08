package shortcuts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Server handles the web interface for shortcuts
type Server struct {
	daemon     *Daemon
	clients    map[*websocket.Conn]bool
	broadcast  chan interface{}
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mu         sync.Mutex
}

// Message represents a WebSocket message
type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

// NewServer creates a new shortcuts server
func NewServer() (*Server, error) {
	daemon, err := NewDaemon()
	if err != nil {
		return nil, fmt.Errorf("failed to create daemon: %w", err)
	}

	s := &Server{
		daemon:     daemon,
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan interface{}),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}

	// Start WebSocket hub
	go s.run()

	return s, nil
}

// Start starts the shortcuts server
func (s *Server) Start() error {
	if err := s.daemon.Start(); err != nil {
		return fmt.Errorf("failed to start daemon: %w", err)
	}

	fmt.Println("🚀 Shortcuts server started")
	return nil
}

// Stop stops the shortcuts server
func (s *Server) Stop() error {
	return s.daemon.Stop()
}

// HandleWebSocket handles WebSocket connections
func (s *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to upgrade connection: %v\n", err)
		return
	}

	s.register <- conn

	defer func() {
		s.unregister <- conn
		conn.Close()
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			fmt.Printf("Failed to parse message: %v\n", err)
			continue
		}

		if err := s.handleMessage(conn, &msg); err != nil {
			fmt.Printf("Failed to handle message: %v\n", err)
		}
	}
}

func (s *Server) run() {
	for {
		select {
		case client := <-s.register:
			s.mu.Lock()
			s.clients[client] = true
			s.mu.Unlock()

		case client := <-s.unregister:
			s.mu.Lock()
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				client.Close()
			}
			s.mu.Unlock()

		case message := <-s.broadcast:
			s.mu.Lock()
			for client := range s.clients {
				if err := client.WriteJSON(message); err != nil {
					fmt.Printf("Failed to send message: %v\n", err)
					client.Close()
					delete(s.clients, client)
				}
			}
			s.mu.Unlock()
		}
	}
}

func (s *Server) handleMessage(conn *websocket.Conn, msg *Message) error {
	switch msg.Type {
	case "get_shortcuts":
		shortcuts, err := s.daemon.GetShortcuts()
		if err != nil {
			return fmt.Errorf("failed to get shortcuts: %w", err)
		}
		return conn.WriteJSON(Message{
			Type:    "shortcuts",
			Payload: shortcuts,
		})

	case "register_shortcut":
		var shortcut Shortcut
		data, err := json.Marshal(msg.Payload)
		if err != nil {
			return fmt.Errorf("failed to marshal shortcut: %w", err)
		}
		if err := json.Unmarshal(data, &shortcut); err != nil {
			return fmt.Errorf("failed to unmarshal shortcut: %w", err)
		}

		if err := s.daemon.RegisterShortcut(shortcut); err != nil {
			return fmt.Errorf("failed to register shortcut: %w", err)
		}

		s.broadcast <- Message{
			Type:    "shortcut_registered",
			Payload: shortcut,
		}

	case "unregister_shortcut":
		var payload struct {
			ID string `json:"id"`
		}
		data, err := json.Marshal(msg.Payload)
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %w", err)
		}
		if err := json.Unmarshal(data, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal payload: %w", err)
		}

		if err := s.daemon.UnregisterShortcut(payload.ID); err != nil {
			return fmt.Errorf("failed to unregister shortcut: %w", err)
		}

		s.broadcast <- Message{
			Type:    "shortcut_unregistered",
			Payload: payload.ID,
		}
	}

	return nil
}
