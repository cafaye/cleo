package release

import "fmt"

func BuildUnknownError(name string) error {
	return fmt.Errorf("unknown release command: %s", name)
}
