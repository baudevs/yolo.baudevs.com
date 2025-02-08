package core

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type GitCommit struct {
	Hash        string
	Author      string
	Date        time.Time
	Message     string
	Files       []string
	Type        string
	Description string
	Impact      string
}

type YoloChange struct {
	Type        string   `yaml:"type"`
	Description string   `yaml:"description"`
	Impact      string   `yaml:"impact"`
	Files       []string `yaml:"files"`
	Status      string   `yaml:"status"`
}

type YoloVersion struct {
	Version string       `yaml:"version"`
	Date    string       `yaml:"date"`
	Changes []YoloChange `yaml:"changes"`
}

func ParseGitHistory() ([]YoloVersion, error) {
	// Check if git repo exists
	if !isGitRepo() {
		return nil, fmt.Errorf("not a git repository")
	}

	commits, err := getGitCommits()
	if err != nil {
		return nil, err
	}

	versions := groupCommitsIntoVersions(commits)
	return versions, nil
}

func isGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	return cmd.Run() == nil
}

func getGitCommits() ([]GitCommit, error) {
	// Get git log with format: hash|author|date|message
	cmd := exec.Command("git", "log", "--pretty=format:%H|%an|%aI|%s")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get git log: %w", err)
	}

	var commits []GitCommit
	lines := strings.Split(string(output), "\n")
	
	for _, line := range lines {
		parts := strings.Split(line, "|")
		if len(parts) != 4 {
			continue
		}

		// Get changed files for this commit
		filesCmd := exec.Command("git", "show", "--name-only", "--format=", parts[0])
		filesOutput, err := filesCmd.Output()
		if err != nil {
			continue
		}
		files := strings.Split(strings.TrimSpace(string(filesOutput)), "\n")

		date, _ := time.Parse(time.RFC3339, parts[2])
		
		commit := GitCommit{
			Hash:   parts[0],
			Author: parts[1],
			Date:   date,
			Message: parts[3],
			Files:  files,
		}

		// Parse conventional commit
		commit.parseConventionalCommit()
		commits = append(commits, commit)
	}

	return commits, nil
}

func (c *GitCommit) parseConventionalCommit() {
	msg := c.Message
	
	// Parse type
	if idx := strings.Index(msg, ":"); idx > 0 {
		c.Type = strings.TrimSpace(msg[:idx])
		msg = strings.TrimSpace(msg[idx+1:])
	} else {
		c.Type = "other"
	}

	// Parse description and impact
	if idx := strings.Index(msg, "("); idx > 0 && strings.Contains(msg, ")") {
		c.Description = strings.TrimSpace(msg[:idx])
		c.Impact = strings.TrimSpace(msg[strings.Index(msg, ")")+1:])
	} else {
		c.Description = msg
		c.Impact = "No impact specified"
	}
}

func groupCommitsIntoVersions(commits []GitCommit) []YoloVersion {
	// Group commits by version tag
	versions := make([]YoloVersion, 0)
	currentVersion := YoloVersion{
		Version: "0.1.0",
		Date:    time.Now().Format("2006-01-02"),
		Changes: make([]YoloChange, 0),
	}

	for _, commit := range commits {
		change := YoloChange{
			Type:        commit.Type,
			Description: commit.Description,
			Impact:      commit.Impact,
			Files:       commit.Files,
			Status:      "implemented",
		}
		currentVersion.Changes = append(currentVersion.Changes, change)
	}

	versions = append(versions, currentVersion)
	return versions
}

func (v YoloVersion) ToYAML() ([]byte, error) {
	return yaml.Marshal(v)
} 