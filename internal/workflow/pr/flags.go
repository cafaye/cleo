package pr

import (
	"fmt"
	"strconv"
	"strings"
)

func hasFlag(args []string, flag string) bool {
	for _, a := range args {
		if a == flag {
			return true
		}
	}
	return false
}

func flagValue(args []string, key string) string {
	for i := 0; i < len(args); i++ {
		if args[i] == key && i+1 < len(args) {
			return args[i+1]
		}
	}
	return ""
}

func flagValues(args []string, key string) []string {
	vals := []string{}
	for i := 0; i < len(args); i++ {
		if args[i] == key && i+1 < len(args) {
			vals = append(vals, args[i+1])
		}
	}
	return vals
}

func parseFrom(args []string) (int, error) {
	for i := 0; i < len(args); i++ {
		if args[i] == "--from" && i+1 < len(args) {
			n, err := strconv.Atoi(args[i+1])
			if err != nil {
				return 0, fmt.Errorf("invalid --from value: %s", args[i+1])
			}
			return n, nil
		}
	}
	return 0, nil
}

func requireLen(args []string, min int, usage string) error {
	if len(args) < min {
		return fmt.Errorf("usage: %s", usage)
	}
	return nil
}

func requireBase(args []string) (string, error) {
	base := strings.TrimSpace(flagValue(args, "--base"))
	if base == "" {
		return "", fmt.Errorf("--base is required")
	}
	return base, nil
}
