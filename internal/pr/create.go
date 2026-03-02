package pr

import (
	"fmt"
	"os"
	"strings"
)

func (s *Service) Create(title, summary, why, what, test, risk, rollback, owner string, cmds []string, draft bool) error {
	branch, err := runLocal("git", "branch", "--show-current")
	if err != nil {
		return err
	}
	head := strings.TrimSpace(branch)
	if head == "" {
		return fmt.Errorf("cannot determine current branch")
	}
	if title == "" {
		title = summary
	}
	if title == "" {
		return fmt.Errorf("title or summary is required")
	}
	body := Render(summary, why, what, test, risk, rollback, owner, cmds)
	tmp, err := os.CreateTemp("", "cleo-pr-body-*.md")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.WriteString(body); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	args := []string{"pr", "create", "--repo", s.repo(), "--base", s.cfg.GitHub.BaseBranch, "--head", head, "--title", title, "--body-file", tmp.Name()}
	if draft {
		args = append(args, "--draft")
	}
	_, err = s.gh.Run(args...)
	return err
}

func Render(summary, why, what, test, risk, rollback, owner string, cmds []string) string {
	if summary == "" {
		summary = "TBD"
	}
	if why == "" {
		why = "TBD"
	}
	if what == "" {
		what = "- TBD"
	}
	if test == "" {
		test = "- TBD"
	}
	if risk == "" {
		risk = "Low"
	}
	if rollback == "" {
		rollback = "Revert this PR"
	}
	if owner == "" {
		owner = "TBD"
	}
	lines := []string{"- `None`"}
	if len(cmds) > 0 {
		lines = []string{}
		for _, c := range cmds {
			lines = append(lines, "- `"+c+"`")
		}
	}
	return fmt.Sprintf(`## Summary
%s

## Why
%s

## What Changed
%s

## How To Test
%s

## Risk
%s

## Rollback
%s

## Ownership
- Primary: %s
- Backup: TBD

## Observability
- Expected signals: TBD
- Dashboard/alerts: TBD

## Post-Merge Production Commands
<!-- post-merge-commands:start -->
%s
<!-- post-merge-commands:end -->
`, summary, why, what, test, risk, rollback, owner, strings.Join(lines, "\n"))
}
