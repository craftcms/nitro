package action

// Destroy will destroy a machine, with an option to permanently delete it.
func Destroy(name string, force bool) (*Action, error) {
	if force {
		return &Action{
			Type:       "delete",
			UseSyscall: false,
			Args:       []string{"delete", name, "-p"},
		}, nil
	}

	return &Action{
		Type:       "delete",
		UseSyscall: false,
		Args:       []string{"delete", name},
	}, nil
}
