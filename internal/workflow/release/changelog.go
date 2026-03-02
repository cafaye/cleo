package release

import (
	"fmt"
	"os"
	"strings"
)

type ChangelogSections struct {
	Summary         string
	Highlights      string
	BreakingChanges string
	MigrationNotes  string
	Verification    string
}

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
	return entry, nil
}

func changelogSections(path, version string) (ChangelogSections, error) {
	entry, err := changelogEntry(path, version)
	if err != nil {
		return ChangelogSections{}, err
	}
	sections := ChangelogSections{
		Summary:         sectionBody(entry, "Summary"),
		Highlights:      sectionBody(entry, "Highlights"),
		BreakingChanges: sectionBody(entry, "Breaking Changes"),
		MigrationNotes:  sectionBody(entry, "Migration Notes"),
		Verification:    sectionBody(entry, "Verification"),
	}
	if err := validateChangelogSections(sections, version); err != nil {
		return ChangelogSections{}, err
	}
	return sections, nil
}

func validateChangelogSections(s ChangelogSections, version string) error {
	missing := []string{}
	if strings.TrimSpace(s.Summary) == "" {
		missing = append(missing, "### Summary")
	}
	if strings.TrimSpace(s.Highlights) == "" {
		missing = append(missing, "### Highlights")
	}
	if strings.TrimSpace(s.BreakingChanges) == "" {
		missing = append(missing, "### Breaking Changes")
	}
	if strings.TrimSpace(s.MigrationNotes) == "" {
		missing = append(missing, "### Migration Notes")
	}
	if strings.TrimSpace(s.Verification) == "" {
		missing = append(missing, "### Verification")
	}
	if len(missing) > 0 {
		return fmt.Errorf("changelog entry for %s is missing required sections: %s", version, strings.Join(missing, ", "))
	}
	if looksPlaceholder(s.Summary) || looksPlaceholder(s.Highlights) || looksPlaceholder(s.BreakingChanges) || looksPlaceholder(s.MigrationNotes) || looksPlaceholder(s.Verification) {
		return fmt.Errorf("changelog entry for %s contains placeholder text", version)
	}
	return nil
}

func sectionBody(entry, heading string) string {
	lines := strings.Split(entry, "\n")
	start := -1
	target := "### " + heading
	for i, line := range lines {
		if strings.TrimSpace(line) == target {
			start = i + 1
			break
		}
	}
	if start == -1 {
		return ""
	}
	end := len(lines)
	for i := start; i < len(lines); i++ {
		if strings.HasPrefix(strings.TrimSpace(lines[i]), "### ") {
			end = i
			break
		}
	}
	return strings.TrimSpace(strings.Join(lines[start:end], "\n"))
}

func looksPlaceholder(text string) bool {
	lc := strings.ToLower(text)
	bad := []string{
		"tbd",
		"initial release for this version",
		"see github changes for merged pr details",
		"lorem ipsum",
	}
	for _, s := range bad {
		if strings.Contains(lc, s) {
			return true
		}
	}
	return false
}
