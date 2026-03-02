# cleo

Deterministic CLI for GitHub PR operations.

## Setup

1. Install GitHub CLI and authenticate (`gh auth status`).
2. Add `cleo.yml` at repo root.
3. Build:

```bash
go build ./cmd/cleo
```

## Commands

```bash
cleo pr status <pr>
cleo pr gate <pr>
cleo pr checks <pr>
cleo pr watch <pr|sha>
cleo pr run <pr> [--dry]
cleo pr create [--title ...] [--summary ...] [--why ...] [--what ...] [--test ...] [--risk ...] [--rollback ...] [--owner ...] [--cmd ...] [--draft]
cleo pr merge <pr> [--no-watch] [--no-run] [--no-rebase] [--delete-branch]
cleo pr rebase <pr>
cleo pr retarget <pr> --base <branch>
cleo pr batch [--from <pr>] [--no-watch] [--no-run] [--no-rebase]
```

## Tests

```bash
go test ./...
```
