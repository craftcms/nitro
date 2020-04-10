package action

func SSH(name string) (*Action, error) {
	return &Action{
		Type:       "shell",
		UseSyscall: true,
		Args:       []string{"shell", name},
	}, nil
}
