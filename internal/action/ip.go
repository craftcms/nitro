package action

import (
	"encoding/json"
	"os/exec"

	"github.com/craftcms/nitro/internal/cmd"
)

func IP(name string, r cmd.ShellRunner) string {
	execCommand := exec.Command(r.Path(), "list", "--format", "json")

	type listOutput struct {
		List []struct {
			Ipv4    []string `json:"ipv4"`
			Name    string   `json:"name"`
			Release string   `json:"release"`
			State   string   `json:"state"`
		} `json:"list"`
	}

	output := listOutput{}

	out, err := execCommand.CombinedOutput()
	if err != nil {
		return ""
	}

	if err := json.Unmarshal(out, &output); err != nil {
		return ""
	}

	ip := ""
	for _, m := range output.List {
		if m.Name == name && len(m.Ipv4) > 0 {
			ip = m.Ipv4[0]
		}
	}

	return ip
}
