package release

import (
	"fmt"
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

func versionFromArgs(args []string, prefix string) (string, error) {
	version := strings.TrimSpace(flagValue(args, "--version"))
	if version == "" && len(args) > 0 && !strings.HasPrefix(args[0], "--") {
		version = strings.TrimSpace(args[0])
	}
	if version == "" {
		return "", fmt.Errorf("--version is required")
	}
	if !strings.HasPrefix(version, prefix) {
		return "", fmt.Errorf("version must start with %q", prefix)
	}
	return version, nil
}
