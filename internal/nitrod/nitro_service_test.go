package nitrod

// spyChainRunner is used for services that run multiple
// Run commands such as editing an ini file and then
// restarting the php-fpm service after completion
type spyChainRunner struct {
	Commands []string
	Args     []map[string][]string
}

func (r *spyChainRunner) Run(command string, args []string) ([]byte, error) {
	r.Commands = append(r.Commands, command)
	r.Args = append(r.Args, map[string][]string{command: args})

	return []byte("something"), nil
}

type spyServiceRunner struct {
	Command string
	Args    []string
}

func (r *spyServiceRunner) Run(command string, args []string) ([]byte, error) {
	r.Command = command

	r.Args = args

	return []byte("test"), nil
}
