package command

type ShellRunner interface {
	Path() string
	UseSyscall(t bool)
	SetInput(input string) error
	Run(args []string) error
}
