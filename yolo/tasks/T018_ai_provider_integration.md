# [T018] AI Provider Integration for Commits

## Status: Implemented
Created: 2024-03-XX
Last Updated: 2024-03-XX
Feature: [F009] AI-Assisted Commits

## Description
Implement integration with multiple AI providers for commit message generation, including configuration management and fallback handling.

## Requirements
- Support OpenAI API
- Support Anthropic Claude API
- Support Mistral API
- Handle API key configuration
- Provide setup instructions
- Implement fallback to manual input
- Support custom prompts

## Implementation
```go
type AIProvider struct {
    Name     string
    APIKey   string
    BaseURL  string
    Model    string
    Prompts  map[string]string
}

type CommitAI struct {
    Provider  AIProvider
    Context   string
    Template  string
}

func (c *CommitAI) GenerateCommitMessage(changes string) (string, error) {
    if c.Provider.APIKey == "" {
        return c.handleMissingConfig()
    }

    prompt := c.buildPrompt(changes)
    return c.callAPI(prompt)
}

func (c *CommitAI) handleMissingConfig() (string, error) {
    fmt.Println("🔑 No AI provider configured. You can:")
    fmt.Println("1. Configure OpenAI API:")
    fmt.Println("   Get key: https://platform.openai.com/api-keys")
    fmt.Println("   Set: export OPENAI_API_KEY=your_key")
    // ... instructions for other providers
    return "", errors.New("no AI provider configured")
}
```

## Notes
- 2024-03-XX: Task created
- 2024-03-XX: Added provider integrations
- 2024-03-XX: Enhanced error handling
- 2024-03-XX: Added setup instructions

## Related
- Parent: [F009] AI-Assisted Commits
- Dependencies: None
- Implements: AI provider integration for commit messages

## Technical Notes
- Uses environment variables for API keys
- Supports provider-specific prompts
- Handles rate limiting and errors
- Provides clear setup guidance 