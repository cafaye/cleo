package pr

import (
	"fmt"
	"os/exec"
	"strings"
)

func runLocal(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s %s: %s", name, strings.Join(args, " "), strings.TrimSpace(string(out)))
	}
	return string(out), nil
}

func (s *Service) Rebase(pr string) error {
	v, err := s.Get(pr)
	if err != nil {
		return err
	}
	if _, err := s.gh.Run("pr", "checkout", pr, "--repo", s.repo()); err != nil {
		return err
	}
	if _, err := runLocal("git", "fetch", "origin", v.BaseRefName, v.HeadRefName); err != nil {
		return err
	}
	if _, err := runLocal("git", "rebase", "origin/"+v.BaseRefName); err != nil {
		return err
	}
	args := []string{"push", "origin", v.HeadRefName}
	if s.cfg.PR.Stack.ForceWithLease {
		args = append(args, "--force-with-lease")
	}
	if _, err := runLocal("git", args...); err != nil {
		return err
	}
	return nil
}

func (s *Service) Retarget(pr, base string) error {
	_, err := s.gh.Run("pr", "edit", pr, "--repo", s.repo(), "--base", base)
	if err != nil {
		return err
	}
	return nil
}
