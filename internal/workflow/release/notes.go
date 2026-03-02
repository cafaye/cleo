package release

import (
	"fmt"
	"strings"
)

var requiredNoteSections = []string{
	"## Summary",
	"## Highlights",
	"## Breaking Changes",
	"## Migration Notes",
	"## Verification",
	"## GitHub Changes",
	"## Changelog",
	"## Full Changelog",
}

func buildReleaseNotes(version, generated string) string {
	return buildReleaseNotesWithChangelog(version, generated, "", "")
}

func buildReleaseNotesWithChangelog(version, generated, changelogEntry, changelogURL string) string {
	fullChangelog := fmt.Sprintf("https://github.com/cafaye/cleo/commits/%s", version)
	lines := []string{
		"## Summary",
		changelogEntry,
		"",
		"## Highlights",
		changelogEntry,
		"",
		"## Breaking Changes",
		"- None.",
		"",
		"## Migration Notes",
		"- None.",
		"",
		"## Verification",
		"- Release artifacts uploaded and checksums generated.",
		"",
		"## GitHub Changes",
		strings.TrimSpace(generated),
		"",
		"## Changelog",
		changelogURL,
		"",
		"## Full Changelog",
		fullChangelog,
		"",
	}
	return strings.Join(lines, "\n")
}

func validateReleaseNotes(body string) error {
	for _, section := range requiredNoteSections {
		if !strings.Contains(body, section) {
			return fmt.Errorf("release notes missing required section: %s", section)
		}
	}
	return nil
}
