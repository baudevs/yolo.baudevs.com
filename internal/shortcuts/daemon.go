package shortcuts

import (
	"fmt"
)

// Daemon represents a system-wide keyboard shortcuts daemon
type Daemon struct {
	isRunning bool
}

// Shortcut represents a system-wide keyboard shortcut
type Shortcut struct {
	ID          string
	Keys        []string
	Command     string
	Description string
}

// NewDaemon creates a new shortcuts daemon
func NewDaemon() (*Daemon, error) {
	return &Daemon{}, nil
}

// Start starts the shortcuts daemon
func (d *Daemon) Start() error {
	d.isRunning = true
	fmt.Println("🚧 Shortcuts daemon is a work in progress")
	return nil
}

// Stop stops the shortcuts daemon
func (d *Daemon) Stop() error {
	d.isRunning = false
	return nil
}

// IsRunning checks if the daemon is running
func (d *Daemon) IsRunning() bool {
	return d.isRunning
}

// RegisterShortcut registers a new shortcut
func (d *Daemon) RegisterShortcut(shortcut Shortcut) error {
	fmt.Printf("🚧 Would register shortcut: %s (%s)\n", shortcut.Description, shortcut.Command)
	return nil
}

// UnregisterShortcut removes a registered shortcut
func (d *Daemon) UnregisterShortcut(id string) error {
	fmt.Printf("🚧 Would unregister shortcut: %s\n", id)
	return nil
}

// GetShortcuts returns all registered shortcuts
func (d *Daemon) GetShortcuts() ([]Shortcut, error) {
	return []Shortcut{}, nil
}
