package action

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"
)

// Update will perform system updates on a given machine
func Update(c *cli.Context) error {
	machine := c.String("machine")
	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))

	cmd := exec.Command(multipass, "exec", machine, "--", "sudo", "bash", "/opt/nitro/update.sh")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Start()
}

func RedisCLIShell(c *cli.Context, e CommandLineExecutor) error {
	machine := c.String("machine")

	args := []string{"multipass", "exec", machine, "--", "redis-cli"}

	return e.Exec(e.Path(), args, os.Environ())
}

// SSH will login a user to a specific machine
func SSH(m string, e CommandLineExecutor) error {
	return e.Exec(e.Path(), []string{"multipass", "shell", m}, os.Environ())
}

func Delete(c *cli.Context) error {
	machine := c.String("machine")

	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))

	cmd := exec.Command(multipass, "delete", machine)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func IP(c *cli.Context) error {
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

func Stop(c *cli.Context) error {
	machine := c.String("machine")

	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))

	cmd := exec.Command(multipass, "stop", machine)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func DatabasePassword(c *cli.Context, e CommandLineExecutor) error {
	return e.Exec(e.Path(), []string{"multipass", "exec", c.String("machine"), "--", "cat", "/home/ubuntu/.db_password"}, os.Environ())
}
