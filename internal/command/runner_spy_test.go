package command

type SpyTestRunner struct {
	args    [][]string
	syscall bool
	input   string
}

func (r *SpyTestRunner) Run(args []string) error {
	r.args = append(r.args, args)

	return nil
}

func (r *SpyTestRunner) UseSyscall(t bool) {
	r.syscall = t
}

func (r *SpyTestRunner) SetInput(input string) error {
	r.input = input

	return nil
}

func (r SpyTestRunner) Path() string {
	return "spytestrunner"
}
