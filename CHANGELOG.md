# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Summary

- Release workflow now requires sectioned changelog input for agent-authored notes.

### Highlights

- Added release inspection commands: `cleo release list` and `cleo release latest`.
- Enforced structured release notes sourced from changelog entries.

### Breaking Changes

- None.

### Migration Notes

- None.

### Verification

- `go test ./...`

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
