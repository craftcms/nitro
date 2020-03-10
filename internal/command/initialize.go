package command

import (
	"strconv"
	"strings"

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
			&cli.BoolFlag{
				Name:        "bootstrap",
				Usage:       "Bootstrap the machine with defaults",
				Value:       true,
				DefaultText: "true",
			},
			&cli.StringFlag{
				Name:        "php-version",
				Usage:       "Provide version of PHP",
				Value:       "7.4",
				DefaultText: "7.4",
			},
			&cli.StringFlag{
				Name:        "database",
				Usage:       "Provide version of PHP",
				Value:       "mariadb",
				DefaultText: "mariadb",
			},
			&cli.Int64Flag{
				Name:        "cpus",
				Usage:       "The number of CPUs to assign the machine",
				Required:    false,
				Value:       2,
				DefaultText: "2",
			},
			&cli.StringFlag{
				Name:        "memory",
				Usage:       "The amount of memory to assign the machine",
				Required:    false,
				Value:       "2G",
				DefaultText: "2G",
			},
			&cli.StringFlag{
				Name:        "disk",
				Usage:       "The amount of disk to assign the machine",
				Required:    false,
				Value:       "5G",
				DefaultText: "5G",
			},
		},
	}
}

func initializeAction(c *cli.Context, r internal.Runner) error {
	machine := c.String("machine")
	cpus := strconv.Itoa(c.Int("cpus"))
	disk := c.String("disk")
	mem := c.String("memory")

	// pass the cloud init as stdin
	r.SetReader(strings.NewReader(cloudInit))

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
