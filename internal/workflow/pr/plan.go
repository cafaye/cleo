package pr

import "fmt"

func BuildPlan(in Input) (Plan, error) {
	switch in.Name {
	case "status":
		if err := requireLen(in.Args, 1, "cleo pr status <pr>"); err != nil {
			return Plan{}, err
		}
		return Plan{Name: in.Name, Description: "Show PR summary", ReadOnly: true}, nil
	case "gate":
		if err := requireLen(in.Args, 1, "cleo pr gate <pr>"); err != nil {
			return Plan{}, err
		}
		return Plan{Name: in.Name, Description: "Validate merge readiness", ReadOnly: true}, nil
	case "checks":
		if err := requireLen(in.Args, 1, "cleo pr checks <pr>"); err != nil {
			return Plan{}, err
		}
		return Plan{Name: in.Name, Description: "Show PR checks", ReadOnly: true}, nil
	case "watch":
		if err := requireLen(in.Args, 1, "cleo pr watch <pr|sha>"); err != nil {
			return Plan{}, err
		}
		return Plan{Name: in.Name, Description: "Watch workflow run", ReadOnly: true}, nil
	case "doctor":
		if len(in.Args) > 0 {
			return Plan{}, fmt.Errorf("usage: cleo pr doctor")
		}
		return Plan{Name: in.Name, Description: "Check PR tooling", ReadOnly: true}, nil
	case "run":
		if err := requireLen(in.Args, 1, "cleo pr run <pr> [--dry]"); err != nil {
			return Plan{}, err
		}
		return Plan{Name: in.Name, Description: "Run post-merge commands", ReadOnly: hasFlag(in.Args, "--dry")}, nil
	case "merge":
		if err := requireLen(in.Args, 1, "cleo pr merge <pr> [--no-watch] [--no-run] [--no-rebase] [--delete-branch]"); err != nil {
			return Plan{}, err
		}
		return Plan{Name: in.Name, Description: "Merge PR with configured safeguards"}, nil
	case "batch":
		if _, err := parseFrom(in.Args); err != nil {
			return Plan{}, err
		}
		return Plan{Name: in.Name, Description: "Merge open PRs in sequence"}, nil
	case "rebase":
		if err := requireLen(in.Args, 1, "cleo pr rebase <pr>"); err != nil {
			return Plan{}, err
		}
		return Plan{Name: in.Name, Description: "Rebase PR branch"}, nil
	case "retarget":
		if err := requireLen(in.Args, 3, "cleo pr retarget <pr> --base <branch>"); err != nil {
			return Plan{}, err
		}
		if _, err := requireBase(in.Args[1:]); err != nil {
			return Plan{}, err
		}
		return Plan{Name: in.Name, Description: "Change PR base branch"}, nil
	case "create":
		return Plan{Name: in.Name, Description: "Create PR with cleo template"}, nil
	default:
		return Plan{}, fmt.Errorf("unknown pr command: %s", in.Name)
	}
}
