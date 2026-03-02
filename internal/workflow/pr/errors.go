package pr

import "fmt"

func BuildUnknownError(name string) error {
	return fmt.Errorf("unknown pr command: %s", name)
}
