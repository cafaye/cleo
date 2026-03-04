package qacatalog

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureQAKitCreatesCoreActor(t *testing.T) {
	root := t.TempDir()
	if err := EnsureQAKit(root); err != nil {
		t.Fatalf("EnsureQAKit() error = %v", err)
	}
	path := filepath.Join(root, ".cleo", "qa", "actors", "core.yml")
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read core actor: %v", err)
	}
	if len(body) == 0 {
		t.Fatal("expected non-empty core actor file")
	}
}

func TestEnsureQAKitDoesNotOverwriteCoreActor(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, ".cleo", "qa", "actors", "core.yml")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	custom := []byte("name: core\ndescription: custom\n")
	if err := os.WriteFile(path, custom, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := EnsureQAKit(root); err != nil {
		t.Fatalf("EnsureQAKit() error = %v", err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(custom) {
		t.Fatal("expected existing core actor file to be preserved")
	}
}
