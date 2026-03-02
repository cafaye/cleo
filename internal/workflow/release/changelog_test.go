package release

import (
	"os"
	"path/filepath"
	"testing"
)

func TestChangelogEntryFound(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "CHANGELOG.md")
	content := "# Changelog\n\n## [v1.0.0]\n### Summary\n- Stable workflow.\n### Highlights\n- Added release automation.\n### Breaking Changes\n- None.\n### Migration Notes\n- None.\n### Verification\n- go test ./...\n\n## [Unreleased]\n- Next\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write changelog: %v", err)
	}
	entry, warnings := changelogSections(path, "v1.0.0")
	if len(warnings) > 0 {
		t.Fatalf("unexpected warnings: %v", warnings)
	}
	if entry.Summary == "" {
		t.Fatal("expected entry")
	}
}

func TestChangelogEntryFallbacks(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "CHANGELOG.md")
	content := "# Changelog\n\n## [v1.0.1]\n### Summary\n- done\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write changelog: %v", err)
	}
	entry, warnings := changelogSections(path, "v1.0.1")
	if len(warnings) == 0 {
		t.Fatal("expected warnings for missing sections")
	}
	if entry.Highlights == "" || entry.BreakingChanges == "" || entry.MigrationNotes == "" || entry.Verification == "" {
		t.Fatal("expected defaults for missing sections")
	}
}
