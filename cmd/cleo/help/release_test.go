package help

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintRelease(t *testing.T) {
	var out bytes.Buffer
	PrintRelease(&out)
	testContainsAll(t, out.String(),
		"usage: cleo release <command>",
		"plan --version",
		"publish --version",
		"help [command]",
	)
}

func TestPrintReleaseCommandKnown(t *testing.T) {
	var out bytes.Buffer
	ok := PrintReleaseCommand(&out, "publish")
	if !ok {
		t.Fatal("expected publish command help")
	}
	if !strings.Contains(out.String(), "usage: cleo release publish") {
		t.Fatalf("unexpected output: %q", out.String())
	}
}

func TestPrintReleaseCommandUnknown(t *testing.T) {
	var out bytes.Buffer
	ok := PrintReleaseCommand(&out, "nope")
	if ok {
		t.Fatal("expected unknown command")
	}
}
