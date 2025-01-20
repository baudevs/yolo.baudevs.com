package git

import (
	"fmt"
	"os/exec"
	"strings"
)

type GitOps struct {
	workingDir string
}

func NewGitOps(workingDir string) *GitOps {
	return &GitOps{
		workingDir: workingDir,
	}
}

func (g *GitOps) GetChanges() (string, error) {
	// Get staged changes
	staged, err := g.runGit("diff", "--cached")
	if err != nil {
		return "", fmt.Errorf("failed to get staged changes: %w", err)
	}

	// Get unstaged changes
	unstaged, err := g.runGit("diff")
	if err != nil {
		return "", fmt.Errorf("failed to get unstaged changes: %w", err)
	}

	// Get untracked files
	untracked, err := g.runGit("ls-files", "--others", "--exclude-standard")
	if err != nil {
		return "", fmt.Errorf("failed to get untracked files: %w", err)
	}

	var changes strings.Builder
	if staged != "" {
		changes.WriteString("Staged changes:\n")
		changes.WriteString(staged)
		changes.WriteString("\n")
	}
	if unstaged != "" {
		changes.WriteString("Unstaged changes:\n")
		changes.WriteString(unstaged)
		changes.WriteString("\n")
	}
	if untracked != "" {
		changes.WriteString("Untracked files:\n")
		changes.WriteString(untracked)
		changes.WriteString("\n")
	}

	return changes.String(), nil
}

func (g *GitOps) StageAll() error {
	_, err := g.runGit("add", ".")
	return err
}

func (g *GitOps) Commit(message string) error {
	_, err := g.runGit("commit", "-m", message)
	return err
}

func (g *GitOps) Push() error {
	// Get current branch
	branch, err := g.getCurrentBranch()
	if err != nil {
		return err
	}

	// Push to remote
	_, err = g.runGit("push", "origin", branch)
	return err
}

func (g *GitOps) Pull() error {
	// Get current branch
	branch, err := g.getCurrentBranch()
	if err != nil {
		return err
	}

	// Pull from remote
	_, err = g.runGit("pull", "origin", branch)
	return err
}

func (g *GitOps) getCurrentBranch() (string, error) {
	output, err := g.runGit("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}
	return strings.TrimSpace(output), nil
}

func (g *GitOps) runGit(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = g.workingDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git command failed: %s: %w", string(output), err)
	}

	return string(output), nil
}

func (g *GitOps) HasRemote() bool {
	_, err := g.runGit("remote")
	return err == nil
}

func (g *GitOps) HasChanges() bool {
	changes, err := g.GetChanges()
	if err != nil {
		return false
	}
	return strings.TrimSpace(changes) != ""
}
