package nitro

// Destroy will destroy a machine, with an option to permanently delete it.
func Destroy(name string) (*Action, error) {
	return &Action{
		Type:       "delete",
		UseSyscall: false,
		Args:       []string{"delete", name, "-p"},
	}, nil
}
