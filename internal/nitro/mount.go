package nitro

func Mount(name, folder, domain string) Command {
	return Command{
		Machine:   name,
		Chainable: true,
		Type:      "mount",
		Args:      []string{folder, name + ":/app/sites/" + domain},
	}
}
