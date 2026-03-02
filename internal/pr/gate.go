package pr

import (
	"fmt"
	"strings"
)

func (s *Service) Gate(pr string) error {
	v, err := s.Get(pr)
	if err != nil {
		return err
	}
	if v.State != "OPEN" {
		return fmt.Errorf("PR #%d is not OPEN (state=%s)", v.Number, v.State)
	}
	if s.cfg.PR.RequireNonDraft && v.IsDraft {
		return fmt.Errorf("PR #%d is draft", v.Number)
	}
	if s.cfg.PR.RequireMergeable && v.Mergeable != "MERGEABLE" {
		return fmt.Errorf("PR #%d is not mergeable (mergeable=%s)", v.Number, v.Mergeable)
	}
	if s.cfg.PR.BlockRequestedChanges && strings.EqualFold(v.ReviewDecision, "CHANGES_REQUESTED") {
		return fmt.Errorf("PR #%d has requested changes", v.Number)
	}
	if err := s.checkRollup(v); err != nil {
		return err
	}
	fmt.Printf("PR #%d is gate-ready.\n", v.Number)
	return nil
}

func (s *Service) checkRollup(v *PRView) error {
	bad := []string{}
	for _, c := range v.StatusCheckRollup {
		if ignored(c.Name, s.cfg.PR.Checks.Ignore) {
			continue
		}
		if c.Status != "COMPLETED" {
			bad = append(bad, fmt.Sprintf("%s/%s status=%s", valueOr(c.WorkflowName, "check"), valueOr(c.Name, "unknown"), c.Status))
			continue
		}
		okConclusion := c.Conclusion == "SUCCESS" || (s.cfg.PR.Checks.TreatNeutralAsPass && c.Conclusion == "NEUTRAL")
		if !okConclusion {
			bad = append(bad, fmt.Sprintf("%s/%s conclusion=%s", valueOr(c.WorkflowName, "check"), valueOr(c.Name, "unknown"), c.Conclusion))
		}
	}
	if len(bad) > 0 {
		return fmt.Errorf("PR #%d has non-green checks:\n- %s", v.Number, strings.Join(bad, "\n- "))
	}
	return nil
}

func ignored(name string, values []string) bool {
	for _, v := range values {
		if strings.EqualFold(strings.TrimSpace(v), strings.TrimSpace(name)) {
			return true
		}
	}
	return false
}
