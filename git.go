package gbt

import (
	"fmt"
	"strings"
)

// GetGitDiffFiles returns the list of files that are different between commit1..commit2
func GetGitDiffFiles(commit1, commit2 string) ([]string, error) {
	out, err := RunCmdWithResult(fmt.Sprintf("git diff --name-only %s..%s", commit1, commit2))
	if err != nil {
		return nil, fmt.Errorf("failed to get git diff files: %w", err)
	}
	out = strings.TrimSpace(out)
	return strings.Split(out, "\n"), nil
}
