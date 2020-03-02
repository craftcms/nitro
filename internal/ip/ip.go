package ip

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/urfave/cli/v2"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:        "ip",
		Usage:       "Show machine IP address",
		Action: func(c *cli.Context) error {
			return run(c)
		},
	}
}

func run(c *cli.Context) error {
	machine := c.String("machine")

	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))

	cmd := exec.Command(multipass, "list", "--format", "json")

	type listOutput struct {
		List []struct {
			Ipv4    []string `json:"ipv4"`
			Name    string   `json:"name"`
			Release string   `json:"release"`
			State   string   `json:"state"`
		} `json:"list"`
	}

	output := listOutput{}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	if err := json.Unmarshal(out, &output); err != nil {
		return err
	}

	ip := ""
	for _, m := range output.List {
		if m.Name == machine && len(m.Ipv4) > 0 {
			ip = m.Ipv4[0]
		}
	}

	if ip == "" {
		fmt.Println("Could not find an IP for the machine:", machine)
		return nil
	}

	fmt.Println(
		fmt.Sprintf("IP address for %s is:\n%s", machine, ip),
	)

	return nil
}
