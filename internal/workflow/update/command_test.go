package update

import "testing"

func TestResolveRefDefault(t *testing.T) {
	if got := resolveRef(""); got != defaultRef {
		t.Fatalf("expected %q, got %q", defaultRef, got)
	}
}

func TestResolveRefValue(t *testing.T) {
	if got := resolveRef("v0.1.0"); got != "v0.1.0" {
		t.Fatalf("unexpected ref: %q", got)
	}
}
