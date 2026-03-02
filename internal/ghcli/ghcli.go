package ghcli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type Client struct{}

func New() *Client { return &Client{} }

func (c *Client) Run(args ...string) (string, error) {
	cmd := exec.Command("gh", args...)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(errOut.String())
		if msg == "" {
			msg = err.Error()
		}
		return "", fmt.Errorf("gh %s: %s", strings.Join(args, " "), msg)
	}
	return out.String(), nil
}

func DecodeJSON(payload string, target any) error {
	return json.Unmarshal([]byte(payload), target)
}
