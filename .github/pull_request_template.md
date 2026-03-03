## Summary

## Why

## What Changed

## How To Test

## Risk

## Rollback

## Ownership
- Primary:
- Backup:

## Acceptance Criteria
<!-- cleo-ac:start -->
version: 1
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
      - replace_with_evidence_artifact
<!-- cleo-ac:end -->

## Observability
- Expected signals:
- Dashboard/alerts:

## Post-Merge Production Commands
<!-- post-merge-commands:start -->
- `None`
<!-- post-merge-commands:end -->
