package pr

import "fmt"

func (s *Service) Status(pr string) error {
	v, err := s.Get(pr)
	if err != nil {
		return err
	}
	fmt.Printf("PR #%d: %s\n", v.Number, v.Title)
	fmt.Printf("URL: %s\n", v.URL)
	fmt.Printf("State: %s (draft=%t, mergeable=%s)\n", v.State, v.IsDraft, v.Mergeable)
	fmt.Printf("Review decision: %s\n", valueOr(v.ReviewDecision, "UNKNOWN"))
	fmt.Printf("Base/Head: %s <- %s\n", v.BaseRefName, v.HeadRefName)
	fmt.Printf("Checks: %d\n", len(v.StatusCheckRollup))
	return nil
}

func (s *Service) Checks(pr string) error {
	v, err := s.Get(pr)
	if err != nil {
		return err
	}
	if len(v.StatusCheckRollup) == 0 {
		fmt.Printf("No status checks reported for PR #%d yet.\n", v.Number)
		fmt.Printf("Try `cleo pr watch %d` and retry.\n", v.Number)
		return nil
	}
	e := s.evaluateChecks(v)
	fmt.Printf("Checks summary: pending=%d failed=%d total=%d\n", len(e.pending), len(e.failed), len(v.StatusCheckRollup))
	for _, c := range v.StatusCheckRollup {
		fmt.Printf("- %s status=%s conclusion=%s\n", checkLabel(c), valueOr(c.Status, "UNKNOWN"), valueOr(c.Conclusion, "UNKNOWN"))
		if c.URL != "" {
			fmt.Printf("  %s\n", c.URL)
		}
		fmt.Println("  trace: check-run-id unavailable via current rollup API; use URL for traceability.")
	}
	if len(e.pending) > 0 {
		fmt.Printf("Pending checks detected. Run `cleo pr watch %d`.\n", v.Number)
	}
	return nil
}
