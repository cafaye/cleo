package pr

import "testing"

func TestEvaluateChecksClassifiesPendingAndFailed(t *testing.T) {
	cfg := testConfig()
	svc := NewServiceWithRunner(cfg, newFakeRunner())
	view := &PRView{
		StatusCheckRollup: []Check{
			{Name: "unit", WorkflowName: "CI", Status: "IN_PROGRESS", Conclusion: "", URL: "u1"},
			{Name: "lint", WorkflowName: "CI", Status: "COMPLETED", Conclusion: "FAILURE", URL: "u2"},
			{Name: "build", WorkflowName: "CI", Status: "COMPLETED", Conclusion: "SUCCESS", URL: "u3"},
		},
	}
	e := svc.evaluateChecks(view)
	if len(e.pending) != 1 {
		t.Fatalf("expected 1 pending check, got %d", len(e.pending))
	}
	if len(e.failed) != 1 {
		t.Fatalf("expected 1 failed check, got %d", len(e.failed))
	}
}
