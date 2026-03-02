package main

import "testing"

func TestUpdateRefDefault(t *testing.T) {
	if got := updateRef([]string{}); got != defaultUpdateRef {
		t.Fatalf("expected default ref %q, got %q", defaultUpdateRef, got)
	}
}

func TestUpdateRefFlag(t *testing.T) {
	if got := updateRef([]string{"--ref", "9060afa"}); got != "9060afa" {
		t.Fatalf("expected ref override, got %q", got)
	}
}
