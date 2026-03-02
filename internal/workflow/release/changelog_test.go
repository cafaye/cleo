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
	entry, err := changelogSections(path, "v1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Summary == "" {
		t.Fatal("expected entry")
	}
}

func TestChangelogEntryRejectsPlaceholder(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "CHANGELOG.md")
	content := "# Changelog\n\n## [v1.0.1]\n### Summary\n- TBD\n### Highlights\n- x\n### Breaking Changes\n- none\n### Migration Notes\n- none\n### Verification\n- x\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write changelog: %v", err)
	}
	if _, err := changelogSections(path, "v1.0.1"); err == nil {
		t.Fatal("expected placeholder validation error")
	}
}
