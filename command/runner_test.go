package command

type RunnerSpy struct {
	Calls []string
}

func (r *RunnerSpy) Run(args []string) error {
	r.Calls = args
	return nil
}