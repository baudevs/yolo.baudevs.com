// models/commit_message.go
package models

// CommitMessage represents the structured output from AI
type CommitMessage struct {
    Type      string   `json:"type"`
    Scope     string   `json:"scope,omitempty"`
    Subject   string   `json:"subject"`
    Body      string   `json:"body,omitempty"`
    Breaking  bool     `json:"breaking,omitempty"`
    IssueRefs []string `json:"issue_refs,omitempty"`
    CoAuthors []string `json:"co_authors,omitempty"`
}

type CommitOptions struct {
	NoSync bool
	Force  bool
}