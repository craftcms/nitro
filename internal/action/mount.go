package action

func Mount(name, folder, site string) (*Action, error) {
	// TODO add validation
	return &Action{
		Type:       "mount",
		UseSyscall: false,
		Args:       []string{"mount", folder, name + ":/app/sites/" + site},
	}, nil
}
