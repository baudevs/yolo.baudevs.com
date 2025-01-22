package messages

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// PersonalityLevel represents the AI's personality level
type PersonalityLevel int

const (
	// Unknown is an unknown personality level
	Unknown PersonalityLevel = iota
	// NerdyClean is a nerdy but clean personality
	NerdyClean
	// MildlyRude is a mildly rude personality
	MildlyRude
	// UnhingedFunny is an unhinged but funny personality
	UnhingedFunny
)

var (
	// Current personality level
	currentPersonality = NerdyClean
)

// GetPersonality returns the current personality level
func GetPersonality() PersonalityLevel {
	return currentPersonality
}

// SetPersonality sets the personality level
func SetPersonality(level PersonalityLevel) {
	currentPersonality = level
}

// GetPersonalityFromString converts a string to a personality level
func GetPersonalityFromString(s string) PersonalityLevel {
	switch strings.ToLower(s) {
	case "nerdy", "clean":
		return NerdyClean
	case "rude", "sassy":
		return MildlyRude
	case "unhinged", "funny":
		return UnhingedFunny
	default:
		return Unknown
	}
}

// String returns the string representation of a personality level
func (p PersonalityLevel) String() string {
	switch p {
	case NerdyClean:
		return "nerdy"
	case MildlyRude:
		return "rude"
	case UnhingedFunny:
		return "unhinged"
	default:
		return "unknown"
	}
}

// Message represents a message with variants for each personality level
type Message struct {
	NerdyClean    string `yaml:"nerdy_clean"`
	MildlyRude    string `yaml:"mildly_rude"`
	UnhingedFunny string `yaml:"unhinged_funny"`
}

// DefaultMessages contains the default message catalog
var DefaultMessages = map[string]Message{
	"welcome": {
		NerdyClean:    "🚀 Welcome to YOLO CLI - Your Optimal Life Organizer!",
		MildlyRude:    "🤘 Sup nerd! Welcome to YOLO - Let's break some stuff!",
		UnhingedFunny: "🔥 YOLO CLI in da house! Time to code like nobody's watching!",
	},
	"install_start": {
		NerdyClean:    "Initiating installation sequence with quantum precision...",
		MildlyRude:    "Alright, let's get this party started! Installing stuff...",
		UnhingedFunny: "Hold onto your bits! We're about to go full send on this install!",
	},
	"install_go": {
		NerdyClean:    "Installing Go - The language of gophers 🐹",
		MildlyRude:    "Yo, we need Go! Don't ask questions, just let it happen...",
		UnhingedFunny: "Time to inject some Go juice into your machine! YEET! 🚀",
	},
	"install_git": {
		NerdyClean:    "Installing Git - Version control for the win!",
		MildlyRude:    "Need Git because apparently you live in 2025... Installing!",
		UnhingedFunny: "Installing Git cuz we ain't savages! Time to get version controlled! 🎮",
	},
	"install_done": {
		NerdyClean:    "Installation complete! Your development environment has been optimized.",
		MildlyRude:    "Done! Try not to break anything important... or do, I'm not your boss!",
		UnhingedFunny: "BOOM! We're in business! Time to write some legendary code, you beautiful disaster!",
	},
	"init_start": {
		NerdyClean:    "Initializing YOLO project structure with mathematical precision...",
		MildlyRude:    "Let's get this party started! Time to YOLO-ify your project...",
		UnhingedFunny: "YOLO MODE ENGAGED! Prepare for project transformation! 🚀",
	},
	"init_done": {
		NerdyClean:    "Project initialized successfully! Ready for optimal productivity.",
		MildlyRude:    "Project's all set up! Don't mess it up... too much.",
		UnhingedFunny: "BOOM! Your project just got YOLO'd! Let the chaos begin! 🎉",
	},
	"commit_start": {
		NerdyClean:    "Analyzing changes with AI precision...",
		MildlyRude:    "Let's see what mess you've made this time...",
		UnhingedFunny: "Time to let the AI judge your code! No pressure! 😈",
	},
	"commit_done": {
		NerdyClean:    "Changes committed successfully! Your code is now immortalized.",
		MildlyRude:    "Alright, your changes are in! Hope you tested them... maybe.",
		UnhingedFunny: "YEET! Your code is now part of history! No takebacks! 🚀",
	},
}

// Messages catalog - will be loaded from config or defaults
var Messages = make(map[string]Message)

func init() {
	// Load custom prompts if they exist, otherwise use defaults
	loadCustomPrompts()
}

// loadCustomPrompts loads custom prompts from the config file
func loadCustomPrompts() {
	home, err := os.UserHomeDir()
	if err != nil {
		Messages = DefaultMessages
		return
	}

	promptsFile := filepath.Join(home, ".config", "yolo", "prompts.yml")
	data, err := os.ReadFile(promptsFile)
	if err != nil {
		Messages = DefaultMessages
		return
	}

	var prompts struct {
		Messages map[string]Message `yaml:"messages"`
	}
	if err := yaml.Unmarshal(data, &prompts); err != nil {
		Messages = DefaultMessages
		return
	}

	// Use custom prompts, falling back to defaults for missing messages
	Messages = DefaultMessages
	for key, msg := range prompts.Messages {
		Messages[key] = msg
	}
}

// Get returns the appropriate message variant based on current personality level
func Get(key string) string {
	msg, ok := Messages[key]
	if !ok {
		return fmt.Sprintf("Message not found: %s", key)
	}

	switch currentPersonality {
	case MildlyRude:
		return msg.MildlyRude
	case UnhingedFunny:
		return msg.UnhingedFunny
	default:
		return msg.NerdyClean
	}
}

// SelectPersonality prompts the user to select a personality level
func SelectPersonality() PersonalityLevel {
	fmt.Println("Select YOLO's personality level:")
	fmt.Println("1) Clean & Nerdy (Safe for work, still fun)")
	fmt.Println("2) Mildly Eccentric (Slightly edgy, occasional sass)")
	fmt.Println("3) Unhinged & Funny (Full chaos mode, not for the faint of heart)")

	var choice string
	fmt.Print("Enter your choice (1-3) [default: 1]: ")
	fmt.Scanln(&choice)

	switch strings.TrimSpace(choice) {
	case "2":
		return MildlyRude
	case "3":
		return UnhingedFunny
	default:
		return NerdyClean
	}
}
