package git

import (
	"fmt"
	"os/exec"
	"strconv"
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
	// Get file statistics
	stats, err := g.runGit("diff", "--numstat")
	if err != nil {
		return "", fmt.Errorf("failed to get change statistics: %w", err)
	}

	// Get staged file statistics
	stagedStats, err := g.runGit("diff", "--cached", "--numstat")
	if err != nil {
		return "", fmt.Errorf("failed to get staged change statistics: %w", err)
	}

	// Get untracked files (just names)
	untracked, err := g.runGit("ls-files", "--others", "--exclude-standard")
	if err != nil {
		return "", fmt.Errorf("failed to get untracked files: %w", err)
	}

	// Get a smart diff for context
	shortDiff, err := g.getSmartDiff()
	if err != nil {
		return "", fmt.Errorf("failed to get smart diff: %w", err)
	}

	var changes strings.Builder

	// Add statistics summary
	if stagedStats != "" {
		changes.WriteString("Staged changes (lines +/-/total):\n")
		changes.WriteString(g.formatStats(stagedStats))
		changes.WriteString("\n")
	}
	if stats != "" {
		changes.WriteString("Unstaged changes (lines +/-/total):\n")
		changes.WriteString(g.formatStats(stats))
		changes.WriteString("\n")
	}

	// Limit untracked files
	if untracked != "" {
		files := strings.Split(untracked, "\n")
		if len(files) > 5 {
			changes.WriteString(fmt.Sprintf("Untracked files (%d total):\n", len(files)))
			for i := 0; i < 5; i++ {
				if files[i] != "" {
					changes.WriteString("  " + files[i] + "\n")
				}
			}
			changes.WriteString("  ... and more\n")
		} else {
			changes.WriteString("Untracked files:\n")
			changes.WriteString(untracked)
			changes.WriteString("\n")
		}
	}

	if shortDiff != "" {
		changes.WriteString("\nKey changes preview:\n")
		changes.WriteString(shortDiff)
	}

	// Ensure the total output doesn't exceed max length
	result := changes.String()
	const maxLength = 4000 // Conservative limit for token size
	if len(result) > maxLength {
		lines := strings.Split(result, "\n")
		var truncated strings.Builder
		currentLength := 0
		
		// Always include stats
		statsEndIndex := 0
		for i, line := range lines {
			if strings.HasPrefix(line, "Key changes preview:") {
				statsEndIndex = i
				break
			}
			truncated.WriteString(line)
			truncated.WriteString("\n")
			currentLength += len(line) + 1
		}

		if currentLength < maxLength {
			truncated.WriteString("\nKey changes preview (truncated due to size):\n")
			for i := statsEndIndex + 1; i < len(lines); i++ {
				lineLength := len(lines[i]) + 1
				if currentLength+lineLength+50 > maxLength { // Leave room for truncation message
					truncated.WriteString("... (output truncated due to length)\n")
					break
				}
				truncated.WriteString(lines[i])
				truncated.WriteString("\n")
				currentLength += lineLength
			}
		}
		return truncated.String(), nil
	}

	return result, nil
}

// formatStats formats the --numstat output into a more readable format
func (g *GitOps) formatStats(stats string) string {
	if stats == "" {
		return ""
	}

	var formatted strings.Builder
	var totalAdded, totalRemoved int
	lines := strings.Split(stats, "\n")
	fileCount := 0
	const maxFiles = 5

	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 3 {
			added := parts[0]
			removed := parts[1]
			file := parts[2]

			// Parse numbers for total
			if added != "-" {
				if n, err := strconv.Atoi(added); err == nil {
					totalAdded += n
				}
			}
			if removed != "-" {
				if n, err := strconv.Atoi(removed); err == nil {
					totalRemoved += n
				}
			}

			// Only show first few files in detail
			if fileCount < maxFiles {
				if added == "-" {
					added = "bin"
				}
				if removed == "-" {
					removed = "bin"
				}
				formatted.WriteString(fmt.Sprintf("  %s: +%s/-%s\n", file, added, removed))
			}
			fileCount++
		}
	}

	// Add summary line
	if fileCount > maxFiles {
		formatted.WriteString(fmt.Sprintf("  ... and %d more files\n", fileCount-maxFiles))
	}
	formatted.WriteString(fmt.Sprintf("  Total: +%d/-%d across %d files\n", totalAdded, totalRemoved, fileCount))

	return formatted.String()
}

// getSmartDiff returns a smart diff that focuses on the most important changes
func (g *GitOps) getSmartDiff() (string, error) {
	// Get the full diff
	fullDiff, err := g.runGit("diff", "--unified=1")  // Reduced context to 1 line
	if err != nil {
		return "", err
	}

	// Split into file chunks
	chunks := strings.Split(fullDiff, "diff --git")
	if len(chunks) <= 1 {
		return fullDiff, nil // Return as is if it's small
	}

	// Process each file's diff
	var smartDiff strings.Builder
	const maxLinesPerFile = 10  // Reduced from 20
	const maxFiles = 3          // Reduced from 5

	processedFiles := 0
	for _, chunk := range chunks[1:] { // Skip first empty chunk
		if processedFiles >= maxFiles {
			smartDiff.WriteString(fmt.Sprintf("\n... (%d more files changed)\n", len(chunks)-maxFiles-1))
			break
		}

		// Split into lines
		lines := strings.Split(chunk, "\n")
		if len(lines) < 3 {
			continue
		}

		// Get file header (simplified)
		parts := strings.Fields(lines[0])
		if len(parts) > 0 {
			smartDiff.WriteString("=== " + parts[len(parts)-1] + " ===\n")
		}

		// Find the actual diff content (skip headers)
		diffStart := 0
		for i, line := range lines {
			if strings.HasPrefix(line, "@@") {
				diffStart = i
				break
			}
		}

		// Process the actual diff content
		diffLines := lines[diffStart:]
		if len(diffLines) > maxLinesPerFile {
			// Keep the first few lines
			smartDiff.WriteString(strings.Join(diffLines[:maxLinesPerFile/2], "\n"))
			smartDiff.WriteString(fmt.Sprintf("\n... (%d lines skipped) ...\n", len(diffLines)-maxLinesPerFile))
			// And the last few lines
			smartDiff.WriteString(strings.Join(diffLines[len(diffLines)-maxLinesPerFile/2:], "\n"))
		} else {
			smartDiff.WriteString(strings.Join(diffLines, "\n"))
		}
		smartDiff.WriteString("\n")

		processedFiles++
	}

	return smartDiff.String(), nil
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
