package update

import (
	"fmt"
	"os"
	"os/exec"
)

const defaultRef = "master"

type Command struct{}

func New() *Command {
	return &Command{}
}

func (c *Command) Execute(nonInteractive bool, ref string) error {
	ref = resolveRef(ref)
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

func resolveRef(ref string) string {
	if ref == "" {
		return defaultRef
	}
	return ref
}
