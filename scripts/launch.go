package scripts

func Launch(machine, cpus, memory, disk string) []string {
	return []string{"launch", "--name", machine, "--cpus", cpus, "--mem", memory, "--disk", disk, "--cloud-init", "-"}
}
