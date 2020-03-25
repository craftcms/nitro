package command

type ShellRunner interface {
	Run(args []string) error
}
