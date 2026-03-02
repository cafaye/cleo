package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const defaultUpdateRef = "master"

func runUpdate(args []string) error {
	ref := updateRef(args)
	nonInteractive := hasFlag(args, "--non-interactive")
	scriptURL := fmt.Sprintf("https://raw.githubusercontent.com/cafaye/cleo/%s/install.sh", ref)
	cmd := exec.Command("bash", "-c", "curl -fsSL "+scriptURL+" | bash")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if nonInteractive {
		cmd.Env = append(os.Environ(), "NON_INTERACTIVE=1")
	}
	return cmd.Run()
}

func updateRef(args []string) string {
	ref := strings.TrimSpace(flagValue(args, "--ref"))
	if ref == "" {
		return defaultUpdateRef
	}
	return ref
}
