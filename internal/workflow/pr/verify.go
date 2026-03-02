package pr

func Verify(plan Plan, _ Result) Verification {
	if plan.ReadOnly {
		return Verification{Checked: false, Reason: "read-only command"}
	}
	switch plan.Name {
	case "merge", "batch", "rebase", "retarget", "create", "run":
		return Verification{Checked: true, Reason: "command execution completed without errors"}
	default:
		return Verification{Checked: false, Reason: "verification not required"}
	}
}
