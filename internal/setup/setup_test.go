package setup

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultConfigIncludesOwnerAndRepo(t *testing.T) {
	cfg := defaultConfig("cafaye", "cleo")
	if !strings.Contains(cfg, "owner: cafaye") {
		t.Fatal("owner missing")
	}
	if !strings.Contains(cfg, "repo: cleo") {
		t.Fatal("repo missing")
	}
	if !strings.Contains(cfg, "block_if_requested_changes") {
		t.Fatal("expected PR policy key")
	}
}

func TestPathContains(t *testing.T) {
	original := os.Getenv("PATH")
	t.Cleanup(func() { _ = os.Setenv("PATH", original) })
	if err := os.Setenv("PATH", "/usr/bin:/tmp/cleo/bin"); err != nil {
		t.Fatal(err)
	}
	if !pathContains("/tmp/cleo/bin") {
		t.Fatal("expected PATH match")
	}
	if pathContains("/tmp/missing/bin") {
		t.Fatal("did not expect PATH match")
	}
}

func TestCopyExecutableCopiesFile(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src")
	dst := filepath.Join(dir, "dst")
	if err := os.WriteFile(src, []byte("cleo-binary"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := copyExecutable(src, dst); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "cleo-binary" {
		t.Fatalf("unexpected copied data: %s", string(body))
	}
}
