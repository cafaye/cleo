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
	if len(v.StatusCheckRollup) == 0 {
		return fmt.Errorf("PR #%d has no status checks reported yet. Re-run `cleo pr checks %d` and `cleo pr watch %d`", v.Number, v.Number, v.Number)
	}
	e := s.evaluateChecks(v)
	if len(e.pending) > 0 {
		return fmt.Errorf("PR #%d has pending checks:\n- %s\nRun `cleo pr watch %d` and retry gate.", v.Number, strings.Join(e.pending, "\n- "), v.Number)
	}
	if len(e.failed) > 0 {
		return fmt.Errorf("PR #%d has non-green checks:\n- %s", v.Number, strings.Join(e.failed, "\n- "))
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
