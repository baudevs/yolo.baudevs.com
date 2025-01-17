package shortcuts

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/baudevs/yolo-cli/internal/shortcuts/macos"
)

// Shortcut represents a global keyboard shortcut
type Shortcut struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Keys        []string `json:"keys"`
	Command     string   `json:"command"`
	Args        []string `json:"args,omitempty"`
	Description string   `json:"description,omitempty"`
	Enabled     bool     `json:"enabled"`
}

// Daemon manages shortcuts
type Daemon struct {
	shortcuts    map[string]*Shortcut
	configPath   string
	mu           sync.RWMutex
	onShortcut   func(shortcut *Shortcut)
	platformDaemon interface{} // Platform-specific daemon
}

// NewDaemon creates a new shortcuts daemon
func NewDaemon(configPath string) (*Daemon, error) {
	d := &Daemon{
		shortcuts:  make(map[string]*Shortcut),
		configPath: configPath,
	}

	// Initialize platform-specific daemon
	if runtime.GOOS == "darwin" {
		macDaemon, err := macos.NewDaemon()
		if err != nil {
			return nil, fmt.Errorf("failed to create macOS daemon: %w", err)
		}
		d.platformDaemon = macDaemon
	} else {
		return nil, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	if err := d.loadConfig(); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return d, nil
}

// Start initializes the daemon
func (d *Daemon) Start() error {
	d.mu.RLock()
	defer d.mu.RUnlock()

	// Register all enabled shortcuts
	for _, s := range d.shortcuts {
		if s.Enabled {
			if err := d.registerShortcut(s); err != nil {
				return fmt.Errorf("failed to register shortcut %s: %w", s.Name, err)
			}
			fmt.Printf("Registered shortcut: %s (%s)\n", s.Name, s.Keys)
		}
	}

	// Start platform daemon
	if runtime.GOOS == "darwin" {
		if err := d.platformDaemon.(*macos.Daemon).Start(); err != nil {
			return fmt.Errorf("failed to start macOS daemon: %w", err)
		}
	}

	return nil
}

// Stop cleans up resources
func (d *Daemon) Stop() {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Stop platform daemon
	if runtime.GOOS == "darwin" {
		d.platformDaemon.(*macos.Daemon).Stop()
	}

	fmt.Println("Stopping shortcuts daemon")
}

// AddShortcut adds a new shortcut
func (d *Daemon) AddShortcut(s *Shortcut) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.shortcuts[s.ID] = s
	if s.Enabled {
		if err := d.registerShortcut(s); err != nil {
			return fmt.Errorf("failed to register shortcut: %w", err)
		}
		fmt.Printf("Added shortcut: %s (%s)\n", s.Name, s.Keys)
	}

	return d.saveConfig()
}

// RemoveShortcut removes a shortcut
func (d *Daemon) RemoveShortcut(id string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if s, exists := d.shortcuts[id]; exists {
		if s.Enabled {
			if err := d.unregisterShortcut(s); err != nil {
				return fmt.Errorf("failed to unregister shortcut: %w", err)
			}
		}
		fmt.Printf("Removed shortcut: %s\n", s.Name)
	}

	delete(d.shortcuts, id)
	return d.saveConfig()
}

// SetEnabled enables or disables a shortcut
func (d *Daemon) SetEnabled(id string, enabled bool) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	s, exists := d.shortcuts[id]
	if !exists {
		return fmt.Errorf("shortcut %s not found", id)
	}

	if s.Enabled != enabled {
		if enabled {
			if err := d.registerShortcut(s); err != nil {
				return fmt.Errorf("failed to register shortcut: %w", err)
			}
		} else {
			if err := d.unregisterShortcut(s); err != nil {
				return fmt.Errorf("failed to unregister shortcut: %w", err)
			}
		}
	}

	s.Enabled = enabled
	fmt.Printf("%s shortcut: %s\n", map[bool]string{true: "Enabled", false: "Disabled"}[enabled], s.Name)

	return d.saveConfig()
}

// OnShortcut sets the callback for when shortcuts are triggered
func (d *Daemon) OnShortcut(callback func(shortcut *Shortcut)) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.onShortcut = callback

	// Set platform callback
	if runtime.GOOS == "darwin" {
		// Create a wrapper function that converts the macos.Shortcut to our Shortcut type
		macosCallback := func(s *macos.Shortcut) {
			if s == nil {
				return
			}
			callback(&Shortcut{
				ID:          s.ID,
				Name:        s.Name,
				Keys:        s.Keys,
				Command:     s.Command,
				Description: s.Description,
				Enabled:     s.Enabled,
			})
		}
		d.platformDaemon.(*macos.Daemon).SetCallback(macosCallback)
	}
}

// GetShortcuts returns all configured shortcuts
func (d *Daemon) GetShortcuts() []*Shortcut {
	d.mu.RLock()
	defer d.mu.RUnlock()

	shortcuts := make([]*Shortcut, 0, len(d.shortcuts))
	for _, s := range d.shortcuts {
		shortcuts = append(shortcuts, s)
	}
	return shortcuts
}

// Helper functions for platform-specific operations
func (d *Daemon) registerShortcut(s *Shortcut) error {
	if runtime.GOOS == "darwin" {
		return d.platformDaemon.(*macos.Daemon).RegisterShortcut(&macos.Shortcut{
			ID:          s.ID,
			Name:        s.Name,
			Keys:        s.Keys,
			Command:     s.Command,
			Description: s.Description,
			Enabled:     s.Enabled,
		})
	}
	return nil
}

func (d *Daemon) unregisterShortcut(s *Shortcut) error {
	if runtime.GOOS == "darwin" {
		return d.platformDaemon.(*macos.Daemon).UnregisterShortcut(s.ID)
	}
	return nil
}

// Configuration persistence
func (d *Daemon) loadConfig() error {
	data, err := os.ReadFile(d.configPath)
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	var shortcuts []*Shortcut
	if err := json.Unmarshal(data, &shortcuts); err != nil {
		return err
	}

	for _, s := range shortcuts {
		d.shortcuts[s.ID] = s
	}
	return nil
}

func (d *Daemon) saveConfig() error {
	shortcuts := make([]*Shortcut, 0, len(d.shortcuts))
	for _, s := range d.shortcuts {
		shortcuts = append(shortcuts, s)
	}

	data, err := json.MarshalIndent(shortcuts, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(d.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(d.configPath, data, 0644)
} 