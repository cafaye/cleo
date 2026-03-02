# Cleo Architecture

## Goals

- Keep CLI routing thin and stable.
- Keep workflow logic modular and testable.
- Standardize behavior for interactive and agent usage.

## Layering

1. `cmd/cleo`
- CLI entrypoint and help text only.
- Parses top-level command/flags and delegates to workflow modules.

2. `internal/workflow/<name>`
- Workflow command boundary.
- Owns command-level input parsing, plan/run/verify orchestration, and workflow-local types.
- Uses adapters to call lower-level domain services.

3. `internal/<domain>`
- Domain/service logic (for example `internal/pr`, `internal/setup`).
- No CLI parsing logic.

## Plan Run Verify Contract

Each workflow command follows:

1. `plan`
- Validate command shape and required inputs.
- Describe intended action and whether it is read-only.

2. `run`
- Execute the action through workflow adapters/services.
- Return structured result state.

3. `verify`
- Confirm completion criteria for the command.
- Return whether verification was performed and why.

## Current Workflow Modules

- `internal/workflow/pr`: PR command contract with plan/run/verify.
- `internal/workflow/setup`: setup command wrapper.
- `internal/workflow/update`: update command wrapper.

## Extension Rule

New workflow families must be added under `internal/workflow/<name>` and keep `cmd/cleo` free of domain logic.
