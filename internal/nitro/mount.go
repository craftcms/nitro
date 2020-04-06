package nitro

func Mount(name, path, domain string) Command {
	return Command{
		Machine: name,
		Type:    "mount",
		Args:    []string{path, name + ":/app/sites/" + domain},
	}
}
