package command

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// Services is a parent command that allows you to manage a machines services
func Services(r Runner) *cli.Command {
	return &cli.Command{
		Name:  "services",
		Usage: "Start, stop, or restart services on machine",
		Action: func(c *cli.Context) error {
			return cli.ShowSubcommandHelp(c)
		},
		Subcommands: []*cli.Command{
			{
				Name:  "restart",
				Usage: "Restart machine services",
				Action: func(c *cli.Context) error {
					if c.Bool("nginx") {
						return r.Run([]string{"exec", c.String("machine"), "--", "sudo", "service", "nginx", "restart"})
					}

					if c.Bool("mysql") {
						return r.Run([]string{"exec", c.String("machine"), "--", "sudo", "service", "mariadb", "restart"})
					}

					if c.Bool("postgres") {
						return r.Run([]string{"exec", c.String("machine"), "--", "sudo", "service", "postgresql", "restart"})
					}

					if c.Bool("redis") {
						return r.Run([]string{"exec", c.String("machine"), "--", "sudo", "service", "redis-server", "restart"})
					}

					return cli.ShowCommandHelp(c, "restart")
				},
				After: func(c *cli.Context) error {
					if c.Bool("nginx") {
						fmt.Println("restarted nginx service for", c.String("machine"))
					}

					if c.Bool("mysql") {
						fmt.Println("restarted mysql service for", c.String("machine"))
					}

					if c.Bool("postgres") {
						fmt.Println("restarted postgres service for", c.String("machine"))
					}

					if c.Bool("redis") {
						fmt.Println("restarted redis service for", c.String("machine"))
					}

					return nil
				},
				Flags: []cli.Flag{
					serviceMySqlFlag,
					serviceNginxFlag,
					servicePostgresFlag,
					serviceRedisFlag,
				},
			},
			{
				Name:  "stop",
				Usage: "Stop machine services",
				Action: func(c *cli.Context) error {
					if c.Bool("nginx") {
						return r.Run([]string{"exec", c.String("machine"), "--", "sudo", "service", "nginx", "stop"})
					}

					if c.Bool("mysql") {
						return r.Run([]string{"exec", c.String("machine"), "--", "sudo", "service", "mariadb", "stop"})
					}

					if c.Bool("postgres") {
						return r.Run([]string{"exec", c.String("machine"), "--", "sudo", "service", "postgresql", "stop"})
					}

					if c.Bool("redis") {
						return r.Run([]string{"exec", c.String("machine"), "--", "sudo", "service", "redis-server", "stop"})
					}

					return cli.ShowCommandHelp(c, "stop")
				},
				After: func(c *cli.Context) error {
					if c.Bool("nginx") {
						fmt.Println("stopped nginx service for", c.String("machine"))
					}

					if c.Bool("mysql") {
						fmt.Println("stopped mysql service for", c.String("machine"))
					}

					if c.Bool("postgres") {
						fmt.Println("stopped postgres service for", c.String("machine"))
					}

					if c.Bool("redis") {
						fmt.Println("stopped redis service for", c.String("machine"))
					}

					return nil
				},
				Flags: []cli.Flag{
					serviceMySqlFlag,
					serviceNginxFlag,
					servicePostgresFlag,
					serviceRedisFlag,
				},
			},
			{
				Name:  "start",
				Usage: "Start machine services",
				Action: func(c *cli.Context) error {
					if c.Bool("nginx") {
						return r.Run([]string{"exec", c.String("machine"), "--", "sudo", "service", "nginx", "start"})
					}

					if c.Bool("mysql") {
						return r.Run([]string{"exec", c.String("machine"), "--", "sudo", "service", "mariadb", "start"})
					}

					if c.Bool("postgres") {
						return r.Run([]string{"exec", c.String("machine"), "--", "sudo", "service", "postgresql", "start"})
					}

					if c.Bool("redis") {
						return r.Run([]string{"exec", c.String("machine"), "--", "sudo", "service", "redis-server", "start"})
					}

					return cli.ShowCommandHelp(c, "start")
				},
				After: func(c *cli.Context) error {
					if c.Bool("nginx") {
						fmt.Println("started nginx service for", c.String("machine"))
					}

					if c.Bool("mysql") {
						fmt.Println("started mysql service for", c.String("machine"))
					}

					if c.Bool("postgres") {
						fmt.Println("started postgres service for", c.String("machine"))
					}

					if c.Bool("redis") {
						fmt.Println("started redis service for", c.String("machine"))
					}

					return nil
				},
				Flags: []cli.Flag{
					serviceMySqlFlag,
					serviceNginxFlag,
					servicePostgresFlag,
					serviceRedisFlag,
				},
			},
		},
	}
}
