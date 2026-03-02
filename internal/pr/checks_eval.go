package pr

import (
	"fmt"
	"strings"
)

type checksEval struct {
	pending []string
	failed  []string
}

func (s *Service) evaluateChecks(v *PRView) checksEval {
	out := checksEval{pending: []string{}, failed: []string{}}
	for _, c := range v.StatusCheckRollup {
		if ignored(c.Name, s.cfg.PR.Checks.Ignore) {
			continue
		}
		label := checkLabel(c)
		if strings.TrimSpace(c.Status) == "" || strings.TrimSpace(c.Status) != "COMPLETED" {
			out.pending = append(out.pending, fmt.Sprintf("%s status=%s", label, valueOr(c.Status, "UNKNOWN")))
			continue
		}
		okConclusion := c.Conclusion == "SUCCESS" || (s.cfg.PR.Checks.TreatNeutralAsPass && c.Conclusion == "NEUTRAL")
		if !okConclusion {
			out.failed = append(out.failed, fmt.Sprintf("%s conclusion=%s", label, valueOr(c.Conclusion, "UNKNOWN")))
		}
	}
	return out
}

func checkLabel(c Check) string {
	return fmt.Sprintf("%s/%s", valueOr(c.WorkflowName, "check"), valueOr(c.Name, "unknown"))
}
