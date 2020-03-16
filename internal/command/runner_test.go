package command

type TestRunner struct {
	args    []string
	syscall bool
	input   string
}

func (r *TestRunner) Run(args []string) error {
	r.args = args

	return nil
}

func (r TestRunner) UseSyscall(t bool) {
	r.syscall = t
}

func (r TestRunner) SetInput(input string) error {
	r.input = input

	return nil
}

func (r TestRunner) Path() string {
	return ""
}
