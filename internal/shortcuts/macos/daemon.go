package macos

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework Carbon
#import <Cocoa/Cocoa.h>
#import <Carbon/Carbon.h>

// Forward declarations
void startEventTap(void);
void stopEventTap(void);
void registerHotKey(const char* identifier, int keyCode, int modifiers);
void unregisterHotKey(const char* identifier);

// Callback function type
typedef void (*ShortcutCallback)(const char* identifier);

// Set callback
void setShortcutCallback(ShortcutCallback callback);
*/
import "C"
import (
	"fmt"
	"sync"
	"unsafe"
)

// Daemon represents the macOS shortcut daemon
type Daemon struct {
	mu        sync.RWMutex
	shortcuts map[string]*Shortcut
	callback  func(shortcut *Shortcut)
}

// Shortcut represents a global keyboard shortcut
type Shortcut struct {
	ID          string
	Name        string
	Keys        []string
	Command     string
	Description string
	Enabled     bool
}

//export shortcutTriggeredCallback
func shortcutTriggeredCallback(cIdentifier *C.char) {
	identifier := C.GoString(cIdentifier)
	
	// Get daemon instance and find shortcut
	if d := getDaemonInstance(); d != nil {
		d.mu.RLock()
		shortcut, exists := d.shortcuts[identifier]
		callback := d.callback
		d.mu.RUnlock()

		if exists && callback != nil {
			callback(shortcut)
		}
	}
}

var callbackFn = shortcutTriggeredCallback

// NewDaemon creates a new macOS shortcut daemon
func NewDaemon() (*Daemon, error) {
	d := &Daemon{
		shortcuts: make(map[string]*Shortcut),
	}

	// Store daemon instance for callbacks
	setDaemonInstance(d)

	// Set up callback
	C.setShortcutCallback((C.ShortcutCallback)(unsafe.Pointer(&callbackFn)))

	return d, nil
}

// Start initializes the daemon and starts listening for shortcuts
func (d *Daemon) Start() error {
	C.startEventTap()
	return nil
}

// Stop stops the daemon and cleans up resources
func (d *Daemon) Stop() {
	C.stopEventTap()
}

// RegisterShortcut registers a new shortcut
func (d *Daemon) RegisterShortcut(s *Shortcut) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Convert keys to Carbon key codes and modifiers
	keyCode, modifiers, err := ParseKeys(s.Keys)
	if err != nil {
		return fmt.Errorf("failed to parse keys: %w", err)
	}

	// Register with system
	cIdentifier := C.CString(s.ID)
	defer C.free(unsafe.Pointer(cIdentifier))
	C.registerHotKey(cIdentifier, C.int(keyCode), C.int(modifiers))

	// Store shortcut
	d.shortcuts[s.ID] = s
	return nil
}

// UnregisterShortcut removes a registered shortcut
func (d *Daemon) UnregisterShortcut(id string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, exists := d.shortcuts[id]; !exists {
		return fmt.Errorf("shortcut %s not found", id)
	}

	// Unregister from system
	cIdentifier := C.CString(id)
	defer C.free(unsafe.Pointer(cIdentifier))
	C.unregisterHotKey(cIdentifier)

	// Remove from map
	delete(d.shortcuts, id)
	return nil
}

// SetCallback sets the callback function for when shortcuts are triggered
func (d *Daemon) SetCallback(callback func(shortcut *Shortcut)) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.callback = callback
}

// Global instance for C callback
var (
	globalDaemon     *Daemon
	globalDaemonLock sync.RWMutex
)

func setDaemonInstance(d *Daemon) {
	globalDaemonLock.Lock()
	defer globalDaemonLock.Unlock()
	globalDaemon = d
}

func getDaemonInstance() *Daemon {
	globalDaemonLock.RLock()
	defer globalDaemonLock.RUnlock()
	return globalDaemon
}

// Helper function to parse keys into Carbon key codes and modifiers
func parseKeys(keys []string) (keyCode int, modifiers int, err error) {
	// TODO: Implement key parsing
	// Convert keys like "⌘", "⌃", "⌥", "⇧" to Carbon modifiers
	// Convert main key to Carbon key code
	return 0, 0, nil
} 