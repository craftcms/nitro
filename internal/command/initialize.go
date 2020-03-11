package command

import (
	"strconv"

	"github.com/urfave/cli/v2"

	"github.com/craftcms/nitro/internal"
)

// Initialize it the main entry point when calling the init command to start a new machine
func Initialize(r internal.Runner) *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "Initialize new machine",
		Action: func(c *cli.Context) error {
			return initializeAction(c, r)
		},
		After: initializeAfterAction,
		Flags: []cli.Flag{
			bootstrapFlag,
			phpVersionFlag,
			databaseFlag,
			cpusFlag,
			memoryFlag,
			diskFlag,
		},
	}
}

func initializeAction(c *cli.Context, r internal.Runner) error {
	machine := c.String("machine")
	cpus := strconv.Itoa(c.Int("cpus"))
	disk := c.String("disk")
	mem := c.String("memory")

	// pass the cloud init as stdin
	if err := r.SetInput(cloudInit); err != nil {
		return err
	}

	return r.Run([]string{"launch", "--name", machine, "--cpus", cpus, "--disk", disk, "--mem", mem, "--cloud-init", "-"})
}

func initializeAfterAction(c *cli.Context) error {
	// if we are bootstrapping, call the command
	if c.Bool("bootstrap") {
		// we are not passing the flags as they should be in the context already
		return c.App.RunContext(c.Context, []string{c.App.Name, "--machine", c.String("machine"), "bootstrap", "--php-version", c.String("php-version"), "--database", c.String("database")}, )
	}

	return nil
}
