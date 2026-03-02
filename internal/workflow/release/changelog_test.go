package release

import (
	"os"
	"path/filepath"
	"testing"
)

func TestChangelogEntryFound(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "CHANGELOG.md")
	content := "# Changelog\n\n## [v1.0.0]\n- Added stable release workflow.\n\n## [Unreleased]\n- Next\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write changelog: %v", err)
	}
	entry, err := changelogEntry(path, "v1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry == "" {
		t.Fatal("expected entry")
	}
}

func TestChangelogEntryRejectsPlaceholder(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "CHANGELOG.md")
	content := "# Changelog\n\n## [v1.0.1]\n- TBD\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write changelog: %v", err)
	}
	if _, err := changelogEntry(path, "v1.0.1"); err == nil {
		t.Fatal("expected placeholder validation error")
	}
}
