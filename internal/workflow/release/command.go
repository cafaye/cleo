package release

type Command struct {
	actions Actions
	opts    Options
}

func New(actions Actions, opts Options) *Command {
	return &Command{actions: actions, opts: opts}
}

func (c *Command) Execute(name string, args []string) error {
	in := Input{Name: name, Args: args}
	plan, err := BuildPlan(in, c.opts)
	if err != nil {
		return err
	}
	result, err := Execute(c.actions, in, c.opts)
	if err != nil {
		return err
	}
	_ = Verify(plan, result)
	return nil
}
