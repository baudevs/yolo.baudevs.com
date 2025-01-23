package utils

import (
	"fmt"
	"path/filepath"
	"strings"
)

// GenerateID generates a new sequential ID with the given prefix (e.g., "E" for epics, "T" for tasks)
func GenerateID(prefix string) string {
	// Get list of existing files
	files, err := filepath.Glob(filepath.Join("yolo", strings.ToLower(prefix)+"*", "*.md"))
	if err != nil {
		return fmt.Sprintf("%s001", prefix)
	}

	// Find highest number
	maxNum := 0
	for _, file := range files {
		base := filepath.Base(file)
		if len(base) < 4 {
			continue
		}
		numStr := strings.TrimPrefix(strings.TrimSuffix(base, ".md"), prefix)
		num := 0
		fmt.Sscanf(numStr, "%d", &num)
		if num > maxNum {
			maxNum = num
		}
	}

	return fmt.Sprintf("%s%03d", prefix, maxNum+1)
}
