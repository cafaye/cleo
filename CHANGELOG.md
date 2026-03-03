# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Summary

- Document upcoming changes here before the next release.

### Highlights

- Add highlights for unreleased work.

### Breaking Changes

- None.

### Migration Notes

- None.

### Verification

- Add verification commands/results for unreleased work.

## [v0.2.0]

### Summary

- Introduced a reusable, CI-gated QA system for PRs with BDD acceptance criteria and policy-driven execution.

### Highlights

- Added end-to-end QA/task workflows for agent-driven execution and follow-up work tracking.
- Standardized Acceptance Criteria to BDD contract fields (`given`, `when`, `then`) with actor/surface/environment metadata.
- Added QA run modes:
  - `auto` (default): validate automated test coverage against BDD criteria.
  - `manual`: optional exploratory/manual execution path.
  - `pr`: resolve mode from PR QA policy block.
- Added PR QA policy markers and parsing (`cleo-qa-policy:start/end`) to support per-PR mode/workflow settings.
- Added QA reporting integration:
  - `cleo qa report` supports PR publishing when explicitly requested.
  - PR checks now surface QA workflow presence relative to AC/policy.
- Added reusable QA bootstrap:
  - `cleo qa init` installs reusable QA kit files.
  - `cleo setup` now bootstraps QA kit assets on first setup.
- Added CI-gated GitHub QA workflow (`.github/workflows/qa.yml`) that runs QA only after CI success.
- Added deterministic local evidence conventions (`qa.evidence_dir`, default `.cleo/evidence`) and per-session evidence directories.
- Added Playwright Go browser automation setup checks and runtime installation support for QA tooling.

### Breaking Changes

- None.

### Migration Notes

- Repositories should include PR AC markers and BDD-formatted criteria for automated QA execution.
- Repositories can run `cleo qa init` to scaffold QA workflow/template assets where missing.
- Default QA automation flow now favors workflow logs/check results over automatic PR comment/body publishing.

### Verification

- `make quality`
- `cleo qa init`
- `cleo qa scaffold`
- `cleo qa start --source pr --ref <pr> --goals "<goal>"`
- `cleo qa run --session <id> --mode auto`
- `cleo qa finish --session <id> --verdict pass`
- `cleo qa report --session <id>`
- `cleo pr checks <pr>`

## [v0.1.4]

### Summary

- Improved PR check reliability and update UX clarity for agent-driven workflows.

### Highlights

- Improved `cleo pr gate` to block on pending or missing checks with explicit `cleo pr watch <pr>` guidance.
- Improved `cleo pr checks` diagnostics with pending/failed summaries and traceability hints.
- Improved `cleo update` logs with current/latest version visibility and clear progress messaging.
- Added repository agent rule to publish a release after each significant improvement.

### Breaking Changes

- None.

### Migration Notes

- None.

### Verification

- `go test ./...`
- `cleo update`

## [v0.1.3]

### Summary

- Made release packaging reusable across projects without hardcoded cleo-only assumptions.

### Highlights

- Added release config keys for cross-project packaging:
  - `release.binary_name`
  - `release.build_target`
  - `release.changelog_file`
- Generalized Go release artifact naming and build target selection.
- Generalized release scripts to accept configurable binary/build targets.

### Breaking Changes

- None.

### Migration Notes

- Existing `cleo.yml` keeps working with defaults.
- Other projects can set `release.binary_name` and `release.build_target`.

### Verification

- `go test ./...`
- `scripts/release/build-assets.sh v0.1.3 /tmp/cleo-release-test`
- `scripts/release/verify-assets.sh v0.1.3 /tmp/cleo-release-test`

## [v0.1.2]

### Summary

- Improved release authoring ergonomics so agents can always produce structured notes.

### Highlights

- Added release inspection commands: `cleo release list` and `cleo release latest`.
- Added warning-first fallback behavior when changelog note sections are missing.
- Added publish-time release note override flags (`--summary`, `--highlights`, `--breaking`, `--migration`, `--verification`).

### Breaking Changes

- None.

### Migration Notes

- None.

### Verification

- `go test ./...`
- `cleo release help publish`

## [v0.1.1]

### Summary

- Established end-to-end GitHub release automation and release-based updater behavior.

### Highlights

- Added artifact build + checksum verification + publish automation.
- Added explicit Go release path (`cleo release go ...`) and runtime split.
- Improved agent workflow guidance and copy-ready `docs/YOUR_AGENTS.md`.

### Breaking Changes

- None.

### Migration Notes

- None.

### Verification

- Release workflow run succeeded for `v0.1.1`.

## [v0.1.0]

### Summary

- Initial public release of `cleo` workflow-driven PR automation.

### Highlights

- Added PR workflow commands (`status`, `gate`, `checks`, `doctor`, `run`, `create`, `merge`, `rebase`, `retarget`, `batch`).
- Added setup wizard and one-command install/uninstall scripts.
- Introduced modular `plan -> run -> verify` workflow architecture.

### Breaking Changes

- None.

### Migration Notes

- None.

### Verification

- Release artifacts and checksums were published.
