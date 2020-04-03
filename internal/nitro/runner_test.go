package nitro

import "errors"

type SpyRunner struct {
	path         string
	calls        []string
	chainedCalls []string
	input        string
	useSyscall   bool
}

func (r *SpyRunner) Path() string {
	return r.path
}

func (r *SpyRunner) UseSyscall(t bool) {
	r.useSyscall = t
}

func (r *SpyRunner) SetInput(input string) error {
	if input == "" {
		return errors.New("you must provide input")
	}
	r.input = input
	return nil
}

func (r *SpyRunner) Run(args []string) error {
	// if we already passed args, lets add this to the chained
	if r.calls != nil {
		r.chainedCalls = args
		return nil
	}
	r.calls = args
	return nil
}
