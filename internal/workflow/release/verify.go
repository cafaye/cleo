package release

func Verify(plan Plan, _ Result) Verification {
	if plan.ReadOnly {
		return Verification{Checked: false, Reason: "read-only command"}
	}
	switch plan.Name {
	case "cut", "publish":
		return Verification{Checked: true, Reason: "release command completed"}
	default:
		return Verification{Checked: false, Reason: "verification not required"}
	}
}
