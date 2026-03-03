package pr

import "strings"

const (
	acStartMarker = "<!-- cleo-ac:start -->"
	acEndMarker   = "<!-- cleo-ac:end -->"
)

func renderACBlock(ac string) string {
	trimmed := strings.TrimSpace(ac)
	if trimmed == "" {
		trimmed = defaultACScaffold()
	}
	return acStartMarker + "\n" + trimmed + "\n" + acEndMarker
}

func defaultACScaffold() string {
	return `version: 1
name: Acceptance Criteria
criteria:
  - id: c1
    title: Replace with criterion title
    severity: medium
    actors: [core]
    acceptance:
      goal: Replace with behavior goal
      expected_result: Replace with expected result
    execution:
      surface: api
      environment: local
      preconditions:
        example_key: example_value
      steps:
        - action: call_api
          params:
            method: GET
            url: http://localhost:0/placeholder
            output_key: response
    evidence_required:
      - replace_with_evidence_artifact`
}
