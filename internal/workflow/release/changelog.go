package release

import (
	"fmt"
	"os"
	"strings"
)

func changelogEntry(path, version string) (string, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(body), "\n")
	start := -1
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "## ["+version+"]" || trimmed == "## "+version {
			start = i + 1
			break
		}
	}
	if start == -1 {
		return "", fmt.Errorf("missing changelog entry for %s in %s", version, path)
	}
	end := len(lines)
	for i := start; i < len(lines); i++ {
		if strings.HasPrefix(strings.TrimSpace(lines[i]), "## ") {
			end = i
			break
		}
	}
	entry := strings.TrimSpace(strings.Join(lines[start:end], "\n"))
	if entry == "" {
		return "", fmt.Errorf("empty changelog entry for %s in %s", version, path)
	}
	if looksPlaceholder(entry) {
		return "", fmt.Errorf("changelog entry for %s contains placeholder text", version)
	}
	return entry, nil
}

func looksPlaceholder(text string) bool {
	lc := strings.ToLower(text)
	bad := []string{
		"tbd",
		"initial release for this version",
		"see github changes for merged pr details",
		"none.",
	}
	for _, s := range bad {
		if strings.Contains(lc, s) {
			return true
		}
	}
	return false
}
