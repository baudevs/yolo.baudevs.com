package macos

import "fmt"

// Carbon modifier key masks
const (
	cmdKey     = 0x0100
	shiftKey   = 0x0200
	optionKey  = 0x0800
	controlKey = 0x1000
)

// Key code mapping for special keys
var keyCodeMap = map[string]int{
	"A": 0x00, "B": 0x0B, "C": 0x08, "D": 0x02, "E": 0x0E,
	"F": 0x03, "G": 0x05, "H": 0x04, "I": 0x22, "J": 0x26,
	"K": 0x28, "L": 0x25, "M": 0x2E, "N": 0x2D, "O": 0x1F,
	"P": 0x23, "Q": 0x0C, "R": 0x0F, "S": 0x01, "T": 0x11,
	"U": 0x20, "V": 0x09, "W": 0x0D, "X": 0x07, "Y": 0x10,
	"Z": 0x06,

	"1": 0x12, "2": 0x13, "3": 0x14, "4": 0x15, "5": 0x17,
	"6": 0x16, "7": 0x1A, "8": 0x1C, "9": 0x19, "0": 0x1D,

	"Space":     0x31,
	"Return":    0x24,
	"Tab":       0x30,
	"Delete":    0x33,
	"Escape":    0x35,
	"Command":   0x37,
	"Shift":     0x38,
	"CapsLock":  0x39,
	"Option":    0x3A,
	"Control":   0x3B,
	"RightShift":   0x3C,
	"RightOption":  0x3D,
	"RightControl": 0x3E,
	"Function":     0x3F,

	"F1": 0x7A, "F2": 0x78, "F3": 0x63, "F4": 0x76,
	"F5": 0x60, "F6": 0x61, "F7": 0x62, "F8": 0x64,
	"F9": 0x65, "F10": 0x6D, "F11": 0x67, "F12": 0x6F,

	"Left":  0x7B,
	"Right": 0x7C,
	"Down":  0x7D,
	"Up":    0x7E,
}

// Modifier key mapping
var modifierMap = map[string]int{
	"⌘": cmdKey,     // Command
	"⇧": shiftKey,   // Shift
	"⌥": optionKey,  // Option/Alt
	"⌃": controlKey, // Control
}

// ParseKeys converts a slice of key strings into Carbon key code and modifiers
func ParseKeys(keys []string) (keyCode int, modifiers int, err error) {
	if len(keys) == 0 {
		return 0, 0, fmt.Errorf("no keys provided")
	}

	// The last key is the main key, everything else is a modifier
	mainKey := keys[len(keys)-1]
	modifierKeys := keys[:len(keys)-1]

	// Get key code for main key
	var ok bool
	keyCode, ok = keyCodeMap[mainKey]
	if !ok {
		return 0, 0, fmt.Errorf("unknown key: %s", mainKey)
	}

	// Combine modifiers
	for _, mod := range modifierKeys {
		modFlag, ok := modifierMap[mod]
		if !ok {
			return 0, 0, fmt.Errorf("unknown modifier: %s", mod)
		}
		modifiers |= modFlag
	}

	return keyCode, modifiers, nil
} 