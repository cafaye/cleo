package pr

import (
	"strings"
	"testing"
)

func TestRenderIncludesMarkersAndSections(t *testing.T) {
	body := Render("Summary", "Why", "- A", "- B", "Low", "Revert", "alice", []string{"bin/kamal logs"})
	checks := []string{
		"## Summary",
		"## Why",
		"## What Changed",
		"## How To Test",
		"## Post-Merge Production Commands",
		"<!-- post-merge-commands:start -->",
		"<!-- post-merge-commands:end -->",
	}
	for _, c := range checks {
		if !strings.Contains(body, c) {
			t.Fatalf("expected body to contain %q", c)
		}
	}
}

func TestRenderDefaults(t *testing.T) {
	body := Render("", "", "", "", "", "", "", nil)
	if !strings.Contains(body, "TBD") {
		t.Fatal("expected defaults to include TBD")
	}
	if !strings.Contains(body, "- `None`") {
		t.Fatal("expected default None command")
	}
}
