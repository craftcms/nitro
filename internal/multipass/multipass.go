package multipass

type Command struct {
	Args []string
}

func Delete(name string, purge bool) *Command {
	args := []string{"delete", name}
	if purge {
		args = append(args, "--purge")
	}

	return &Command{
		Args: args,
	}
}

func Exec(name string) *Command {
	return &Command{
		Args: []string{"exec", name, "--"},
	}
}

func Info(name, format string) *Command {
	var f string
	switch format {
	case "json":
		f = "json"
	case "csv":
		f = "csv"
	default:
		f = "table"
	}

	return &Command{
		Args: []string{"info", name, "--format", f},
	}
}

func Launch(name, cpus, memory, disk string, cloudConfig bool) *Command {
	args := []string{"--name", name, "--cpus", cpus, "--mem", memory, "--disk", disk}
	if cloudConfig {
		args = append(args, "--cloud-init", "-")
	}

	return &Command{
		Args: args,
	}
}

func Mount(name, path, target string) *Command {
	return &Command{
		Args: []string{"mount", path, name + ":/" + target},
	}
}

func Shell(name string) *Command {
	return &Command{
		Args: []string{"shell", name},
	}
}
