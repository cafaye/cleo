package pr

import "testing"

func TestBuildPlanStatus(t *testing.T) {
	plan, err := BuildPlan(Input{Name: "status", Args: []string{"12"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !plan.ReadOnly {
		t.Fatal("status should be read-only")
	}
}

func TestBuildPlanRetargetNeedsBase(t *testing.T) {
	_, err := BuildPlan(Input{Name: "retarget", Args: []string{"12", "--base"}})
	if err == nil {
		t.Fatal("expected error for missing --base value")
	}
}

func TestBuildPlanUnknown(t *testing.T) {
	_, err := BuildPlan(Input{Name: "nope"})
	if err == nil {
		t.Fatal("expected unknown command error")
	}
}
