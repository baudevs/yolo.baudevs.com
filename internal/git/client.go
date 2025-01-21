package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// Client represents a Git client
type Client struct {
	workingDir string
}

// NewClient creates a new Git client
func NewClient() (*Client, error) {
	// Check if git is installed
	if _, err := exec.LookPath("git"); err != nil {
		return nil, fmt.Errorf("git is not installed: %w", err)
	}

	// Get git root directory
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("not a git repository: %w", err)
	}

	return &Client{
		workingDir: strings.TrimSpace(string(out)),
	}, nil
}

// HasChanges checks if there are any changes in the working directory
func (c *Client) HasChanges() bool {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = c.workingDir

	out, err := cmd.Output()
	if err != nil {
		return false
	}

	return len(out) > 0
}

// GetChanges returns all changes in the working directory
func (c *Client) GetChanges() (string, error) {
	cmd := exec.Command("git", "diff")
	cmd.Dir = c.workingDir

	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get changes: %w", err)
	}

	// Also get untracked files
	cmd = exec.Command("git", "ls-files", "--others", "--exclude-standard")
	cmd.Dir = c.workingDir

	untracked, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get untracked files: %w", err)
	}

	result := string(out)
	if len(untracked) > 0 {
		result += "\n\nUntracked files:\n" + string(untracked)
	}

	return result, nil
}

// StageAll stages all changes
func (c *Client) StageAll() error {
	cmd := exec.Command("git", "add", "-A")
	cmd.Dir = c.workingDir

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stage changes: %w", err)
	}

	return nil
}

// GetStagedDiff returns the diff of staged changes
func (c *Client) GetStagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	cmd.Dir = c.workingDir

	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get staged diff: %w", err)
	}

	return string(out), nil
}

// Commit creates a new commit with the given message
func (c *Client) Commit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Dir = c.workingDir

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}

// HasRemote checks if the repository has a remote
func (c *Client) HasRemote() bool {
	cmd := exec.Command("git", "remote")
	cmd.Dir = c.workingDir

	out, err := cmd.Output()
	if err != nil {
		return false
	}

	return len(out) > 0
}

// Pull pulls changes from remote
func (c *Client) Pull() error {
	cmd := exec.Command("git", "pull", "--rebase")
	cmd.Dir = c.workingDir

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to pull: %w", err)
	}

	return nil
}

// Push pushes changes to remote
func (c *Client) Push() error {
	cmd := exec.Command("git", "push")
	cmd.Dir = c.workingDir

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to push: %w", err)
	}

	return nil
}

// GetLastCommitMessage returns the last commit message
func (c *Client) GetLastCommitMessage() (string, error) {
	cmd := exec.Command("git", "log", "-1", "--pretty=%B")
	cmd.Dir = c.workingDir

	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get last commit message: %w", err)
	}

	return strings.TrimSpace(string(out)), nil
}
