package types

// Config holds configuration
type Config struct {
	APIEndpoint      string            `json:"api_endpoint"`
	OpenAIKey       string            `json:"openai_key"`
	DevMode         bool              `json:"dev_mode"`
	Prompts         map[string]string `json:"prompts"`
	PersonalityType string            `json:"personality_type"`
}
