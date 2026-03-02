package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadAppliesDefaults(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cleo.yml")
	body := []byte("version: 1\ngithub:\n  owner: cafaye\n  repo: cleo\n")
	if err := os.WriteFile(path, body, 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.GitHub.BaseBranch != "main" {
		t.Fatalf("expected default base branch main, got %s", cfg.GitHub.BaseBranch)
	}
	if cfg.PR.Checks.Mode != "required" {
		t.Fatalf("expected default checks mode required, got %s", cfg.PR.Checks.Mode)
	}
}

func TestLoadRequiresRepoFields(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cleo.yml")
	if err := os.WriteFile(path, []byte("version: 1\ngithub:\n  owner: x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil {
		t.Fatal("expected validation error")
	}
}
