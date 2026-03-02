package update

import "fmt"

type Command struct {
	updater *ReleaseUpdater
	current string
}

func New(current string) *Command {
	return &Command{updater: NewReleaseUpdater(), current: current}
}

func (c *Command) Execute(_ bool) error {
	if err := c.updater.UpdateLatest(c.current); err != nil {
		return fmt.Errorf("release update failed: %w", err)
	}
	return nil
}
