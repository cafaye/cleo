package pr

import (
	"fmt"
	"sort"
	"strconv"
)

func (s *Service) Merge(pr string, noWatch bool, noRun bool, noRebase bool, deleteBranch bool) error {
	if err := s.Gate(pr); err != nil {
		return err
	}
	args := []string{"pr", "merge", pr, "--repo", s.repo()}
	switch s.cfg.GitHub.MergeMethod {
	case "squash":
		args = append(args, "--squash")
	case "rebase":
		args = append(args, "--rebase")
	default:
		args = append(args, "--merge")
	}
	if deleteBranch || s.cfg.GitHub.DeleteBranchOnMerge {
		args = append(args, "--delete-branch")
	}
	if _, err := s.gh.Run(args...); err != nil {
		return err
	}
	fmt.Printf("Merged PR #%s.\n", pr)
	v, err := s.Get(pr)
	if err == nil && s.cfg.PR.DeployWatch.Enabled && !noWatch && v.BaseRefName == s.cfg.GitHub.BaseBranch {
		if err := s.Watch(pr); err != nil {
			return err
		}
	}
	if s.cfg.PR.PostMerge.Enabled && !noRun {
		if err := s.Run(pr, false); err != nil {
			return err
		}
	}
	if s.cfg.PR.Stack.RebaseNextAfterMerge && !noRebase {
		if err := s.retargetNext(pr); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) Batch(start int, noWatch bool, noRun bool, noRebase bool) error {
	out, err := s.gh.Run("pr", "list", "--repo", s.repo(), "--state", "open", "--limit", "200", "--json", "number")
	if err != nil {
		return err
	}
	rows, err := ghDecodeRows(out)
	if err != nil {
		return err
	}
	prs := make([]int, 0, len(rows))
	for _, r := range rows {
		if r >= start {
			prs = append(prs, r)
		}
	}
	sort.Ints(prs)
	for _, n := range prs {
		fmt.Printf("-----\nProcessing PR #%d\n", n)
		if err := s.Merge(strconv.Itoa(n), noWatch, noRun, noRebase, false); err != nil {
			return err
		}
	}
	fmt.Println("Done.")
	return nil
}

func (s *Service) retargetNext(currentPR string) error {
	if !s.cfg.PR.Stack.AutoDetectNextPR {
		return nil
	}
	current, err := parseNumber(currentPR)
	if err != nil {
		return nil
	}
	out, err := s.gh.Run("pr", "list", "--repo", s.repo(), "--state", "open", "--limit", "200", "--json", "number")
	if err != nil {
		return err
	}
	rows, err := ghDecodeRows(out)
	if err != nil {
		return err
	}
	next := 0
	for _, n := range rows {
		if n > current && (next == 0 || n < next) {
			next = n
		}
	}
	if next == 0 {
		return nil
	}
	return s.Retarget(strconv.Itoa(next), s.cfg.GitHub.BaseBranch)
}
