# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

- Added release inspection commands: `cleo release list` and `cleo release latest`.
- Enforced structured release notes sourced from changelog entries.

## [v0.1.1]

### Added

- End-to-end GitHub release workflow with artifact build, checksum verification, and publish automation.
- Release-based updater (`cleo update`) using published assets and checksum verification.
- Explicit Go release path (`cleo release go ...`) and Go runtime module split.

### Changed

- Standardized release notes structure and release command output messages.
- Improved agent workflow guidance and copy-ready `docs/YOUR_AGENTS.md` template.

## [v0.1.0]

### Added

- Initial `cleo` PR workflow commands for status, checks, gate, doctor, run, create, merge, rebase, retarget, and batch.
- Setup wizard, one-command install/uninstall scripts, and non-interactive execution support.

### Changed

- Modular workflow architecture with `plan -> run -> verify` command contract.
- Logging defaults and command help coverage for workflow-first usage.
