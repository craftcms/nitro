package scripts

// DockerRunDatabase builds the commands to run a database inside of docker for a machine.
func DockerRunDatabase(name, engine, version string) []string {
	var port string
	var envvars []string
	switch engine {
	case "postgres":
		port = "5432"
		envvars = []string{"-e", "POSTGRES_PASSWORD=nitro", "-e", "POSTGRES_USER=nitro", "-e", "POSTGRES_DB=nitro"}
	default:
		port = "3306"
		envvars = []string{"-e", "MYSQL_ROOT_PASSWORD=nitro", "-e", "MYSQL_DATABASE=nitro", "-e", "MYSQL_USER=nitro", "-e", "MYSQL_PASSWORD=nitro"}
	}

	commands := []string{"exec", name, "--", "docker", "run", "-d", "--restart=always", engine + ":" + version, "-p", port + ":" + port}

	// append the environment variables
	commands = append(commands, envvars...)

	return commands
}
