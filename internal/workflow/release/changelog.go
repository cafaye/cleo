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

func changelogSections(path, version string) (ChangelogSections, []string) {
	entry, err := changelogEntry(path, version)
	if err != nil {
		return defaultSections(), []string{fmt.Sprintf("No changelog entry for %s in %s. Add one or pass release flags.", version, path)}
	}
	sections := ChangelogSections{
		Summary:         sectionBody(entry, "Summary"),
		Highlights:      sectionBody(entry, "Highlights"),
		BreakingChanges: sectionBody(entry, "Breaking Changes"),
		MigrationNotes:  sectionBody(entry, "Migration Notes"),
		Verification:    sectionBody(entry, "Verification"),
	}
	warnings := []string{}
	if strings.TrimSpace(sections.Summary) == "" {
		warnings = append(warnings, "Missing ### Summary in changelog entry.")
		sections.Summary = defaultSections().Summary
	}
	if strings.TrimSpace(sections.Highlights) == "" {
		warnings = append(warnings, "Missing ### Highlights in changelog entry.")
		sections.Highlights = defaultSections().Highlights
	}
	if strings.TrimSpace(sections.BreakingChanges) == "" {
		warnings = append(warnings, "Missing ### Breaking Changes in changelog entry.")
		sections.BreakingChanges = defaultSections().BreakingChanges
	}
	if strings.TrimSpace(sections.MigrationNotes) == "" {
		warnings = append(warnings, "Missing ### Migration Notes in changelog entry.")
		sections.MigrationNotes = defaultSections().MigrationNotes
	}
	if strings.TrimSpace(sections.Verification) == "" {
		warnings = append(warnings, "Missing ### Verification in changelog entry.")
		sections.Verification = defaultSections().Verification
	}
	return sections, warnings
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

func defaultSections() ChangelogSections {
	return ChangelogSections{
		Summary:         "- Provide release summary via CHANGELOG.md `### Summary` or `--summary`.",
		Highlights:      "- Provide highlights via CHANGELOG.md `### Highlights` or `--highlights`.",
		BreakingChanges: "- Provide breaking changes via `### Breaking Changes` or `--breaking` (use `None` if none).",
		MigrationNotes:  "- Provide migration notes via `### Migration Notes` or `--migration` (use `None` if none).",
		Verification:    "- Provide verification notes via `### Verification` or `--verification`.",
	}
}
