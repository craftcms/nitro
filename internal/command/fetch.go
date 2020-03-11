package command

import (
	"encoding/json"
	"errors"
	"os/exec"
)

func FetchIP(machine string, r Runner) (string, error) {
	cmd := exec.Command(r.Path(), "list", "--format", "json")

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
		return "", err
	}

	if err := json.Unmarshal(out, &output); err != nil {
		return "", err
	}

	ip := ""
	for _, m := range output.List {
		if m.Name == machine && len(m.Ipv4) > 0 {
			ip = m.Ipv4[0]
		}
	}

	if ip == "" {
		return "", errors.New("Could not find an IP for: " + machine)
	}

	return ip, nil
}
