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
	// Get summary of changes
	summary, err := g.runGit("diff", "--compact-summary")
	if err != nil {
		return "", fmt.Errorf("failed to get changes summary: %w", err)
	}

	// Get staged changes summary
	stagedSummary, err := g.runGit("diff", "--cached", "--compact-summary")
	if err != nil {
		return "", fmt.Errorf("failed to get staged changes summary: %w", err)
	}

	// Get untracked files (just names)
	untracked, err := g.runGit("ls-files", "--others", "--exclude-standard")
	if err != nil {
		return "", fmt.Errorf("failed to get untracked files: %w", err)
	}

	// Get a short diff for context (limited to first few lines)
	shortDiff, err := g.runGit("diff", "--unified=3")
	if err != nil {
		return "", fmt.Errorf("failed to get short diff: %w", err)
	}

	// Truncate the diff if it's too long
	diffLines := strings.Split(shortDiff, "\n")
	if len(diffLines) > 50 {
		diffLines = diffLines[:50]
		diffLines = append(diffLines, "... (truncated)")
	}
	shortDiff = strings.Join(diffLines, "\n")

	var changes strings.Builder
	if stagedSummary != "" {
		changes.WriteString("Staged changes:\n")
		changes.WriteString(stagedSummary)
		changes.WriteString("\n")
	}
	if summary != "" {
		changes.WriteString("Unstaged changes:\n")
		changes.WriteString(summary)
		changes.WriteString("\n")
	}
	if untracked != "" {
		changes.WriteString("Untracked files:\n")
		changes.WriteString(untracked)
		changes.WriteString("\n")
	}
	if shortDiff != "" {
		changes.WriteString("\nChanges preview:\n")
		changes.WriteString(shortDiff)
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
