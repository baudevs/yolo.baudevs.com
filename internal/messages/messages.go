package messages

import (
	"fmt"
	"os"
	"strings"
)

// PersonalityLevel defines the tone of messages
type PersonalityLevel int

const (
	// NerdyClean - Correct and clean but nerdy funny
	NerdyClean PersonalityLevel = iota + 1
	// MildlyRude - Correct and mildly eccentric with tones of rudeness
	MildlyRude
	// UnhingedFunny - Politically incorrect and funny
	UnhingedFunny
)

var currentLevel PersonalityLevel = NerdyClean

// Message represents a message with variants for each personality level
type Message struct {
	NerdyClean    string
	MildlyRude    string
	UnhingedFunny string
}

// Messages catalog
var Messages = map[string]Message{
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
	"install_success": {
		NerdyClean:    "🎉 Installation complete! Your development environment has been successfully optimized!",
		MildlyRude:    "Done! Try not to break anything important... or do, I'm not your boss! 😏",
		UnhingedFunny: "BOOM! We're in business! Time to write some legendary code, you beautiful disaster! 🚀",
	},
	"init_project": {
		NerdyClean:    "Initializing project structure with mathematical precision...",
		MildlyRude:    "Making folders and stuff. Try to keep them organized this time...",
		UnhingedFunny: "Time to birth a new project! Push! 🤰",
	},
	"commit_start": {
		NerdyClean:    "Analyzing changes with quantum precision...",
		MildlyRude:    "Let's see what mess you've made this time...",
		UnhingedFunny: "Time to immortalize your code crimes! What did you do!? 👮‍♂️",
	},
	"commit_success": {
		NerdyClean:    "Changes committed successfully! Your code is now part of history.",
		MildlyRude:    "Alright, I've hidden your changes in the repo. Happy now?",
		UnhingedFunny: "Your code is now officially in witness protection! 🥸",
	},
	"error_generic": {
		NerdyClean:    "Oops! We've encountered an unexpected quantum fluctuation...",
		MildlyRude:    "Well, that didn't work. Want to try again or...?",
		UnhingedFunny: "FAIL! 💩 Everything's broken! But hey, that's job security!",
	},
	"error_git": {
		NerdyClean:    "Git seems to be having an existential crisis...",
		MildlyRude:    "Git is being Git again. You know how it is...",
		UnhingedFunny: "Git just went YOLO and not in a good way! 🎢",
	},
	"error_ai": {
		NerdyClean:    "The AI is currently contemplating the meaning of life...",
		MildlyRude:    "AI machine broke. Have you tried turning it off and on again?",
		UnhingedFunny: "The AI is having a mental breakdown! Time for therapy! 🛋️",
	},
}

// SetPersonality sets the global personality level
func SetPersonality(level PersonalityLevel) {
	if level < NerdyClean || level > UnhingedFunny {
		level = NerdyClean
	}
	currentLevel = level
	os.Setenv("YOLO_PERSONALITY", fmt.Sprintf("%d", level))
}

// GetPersonality returns the current personality level
func GetPersonality() PersonalityLevel {
	if level := os.Getenv("YOLO_PERSONALITY"); level != "" {
		if l, err := fmt.Sscanf(level, "%d", &currentLevel); err == nil && l > 0 {
			return currentLevel
		}
	}
	return currentLevel
}

// Get returns the appropriate message variant based on current personality level
func Get(key string) string {
	msg, ok := Messages[key]
	if !ok {
		return "Message not found"
	}

	switch GetPersonality() {
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

	choice = strings.TrimSpace(choice)
	switch choice {
	case "2":
		return MildlyRude
	case "3":
		return UnhingedFunny
	default:
		return NerdyClean
	}
}
