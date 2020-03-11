package command

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// IP will look for a specific machine IP address by name
func IP(r Runner) *cli.Command {
	return &cli.Command{
		Name:  "ip",
		Usage: "Show machine IP address",
		Action: func(c *cli.Context) error {
			return ipAction(c, r)
		},
	}
}

func ipAction(c *cli.Context, r Runner) error {
	ip, err := fetchIP(c.String("machine"), r)
	if err != nil {
		return err
	}

	fmt.Println(ip)

	return nil
}
