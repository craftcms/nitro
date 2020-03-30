package x

import "github.com/craftcms/nitro/scripts"

type Command struct {
	Description string
	Machine     string
	Args        map[string][]string
	StdIn       string
}

func Init(machine, cpus, memory, disk, php, dbEngine, dbVersion string) Command {
	// setup the core packages
	launchCommand := scripts.Launch(machine, cpus, memory, disk)
	installCommand := scripts.AptInstall(machine, scripts.InstallPHP(php))
	dockerCommand := scripts.DockerRunDatabase(machine, dbEngine, dbVersion)

	return Command{
		Description: "launch and init",
		Machine:     machine,
		Args: map[string][]string{
			"launch":  launchCommand,
			"install": installCommand,
			"docker":  dockerCommand,
		},
	}
}
