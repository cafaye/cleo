package setup

import "testing"

func TestNew(t *testing.T) {
	if New() == nil {
		t.Fatal("expected command instance")
	}
}
