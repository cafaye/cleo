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
		fmt.Println("No status checks reported.")
		return nil
	}
	for _, c := range v.StatusCheckRollup {
		fmt.Printf("- %s/%s status=%s conclusion=%s\n", valueOr(c.WorkflowName, "check"), valueOr(c.Name, "unknown"), valueOr(c.Status, "UNKNOWN"), valueOr(c.Conclusion, "UNKNOWN"))
		if c.URL != "" {
			fmt.Printf("  %s\n", c.URL)
		}
	}
	return nil
}
