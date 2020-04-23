package find

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os/exec"
	"strings"

	"github.com/craftcms/nitro/config"
)

// Mounts will take a name of a machine and the output of an exec.Command as a slice of bytes
// and return a slice of config mounts that has a source and destination or an error. This is
// used to match if the machine has any mounts. The args passed to multipass are expected to
// return a csv format (e.g. "multipass info machinename --format=csv").
func Mounts(name string, b []byte) ([]config.Mount, error) {
	var mounts []config.Mount

	records, err := csv.NewReader(bytes.NewReader(b)).ReadAll()
	if err != nil {
		return nil, err
	}

	for i, record := range records {
		if i == 0 || record[0] != name {
			continue
		}

		for _, m := range strings.Split(record[12], ";") {
			// since we split on the ;, the last element could be empty
			if len(m) == 0 {
				break
			}

			mount := strings.Split(m, " ")

			mounts = append(mounts, config.Mount{Source: mount[0], Dest: mount[2]})
		}
	}

	return mounts, nil
}

func ContainersToCreate(machine string, cfg config.Config) ([]config.Database, error) {
	path, err := exec.LookPath("multipass")
	if err != nil {
		return nil, err
	}

	var dbs []config.Database
	for _, db := range cfg.Databases {
		container := fmt.Sprintf("%s_%s_%s", db.Engine, db.Version, db.Port)

		c := exec.Command(path, []string{"exec", machine, "--", "sudo", "bash", "/opt/nitro/scripts/docker-container-exists.sh", container}...)
		output, err := c.Output()
		if err != nil {
			return nil, err
		}

		if !strings.Contains(string(output), "exists") {
			dbs = append(dbs, db)
		}
	}

	return dbs, nil
}
