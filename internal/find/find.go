package find

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"os/exec"
	"strings"

	"github.com/craftcms/nitro/config"
)

// Finder is an interface the wraps the exec.Command Output function
// it is used by this package to parse output of the exec.Command
type Finder interface {
	Output() ([]byte, error)
}

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
		c := exec.Command(path, []string{"exec", machine, "--", "sudo", "bash", "/opt/nitro/scripts/docker-container-exists.sh", db.Name()}...)
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

// SitesEnabled takes a finder which is a command executed
// by the multipass cli tool that outputs the contents
// (symlinks) or sites-enabled and returns sites.
func SitesEnabled(f Finder) ([]config.Site, error) {
	out, err := f.Output()
	if err != nil {
		return nil, err
	}

	// parse the out
	var sites []config.Site
	sc := bufio.NewScanner(strings.NewReader(string(out)))
	for sc.Scan() {
		if l := sc.Text(); l != "" {
			sp := strings.Split(strings.TrimSpace(sc.Text()), "/")
			if h := sp[len(sp)-1]; h != "default" {
				sites = append(sites, config.Site{Hostname: h})
			}
		}
	}

	return sites, nil
}

// PHPVersion is used to get the "current" or "default" version
// of PHP that is installed. It expects exec.Command to sent
// "multipass", "exec", machine, "--", "php", "--version"
func PHPVersion(f Finder) (string, error) {
	out, err := f.Output()
	if err != nil {
		return "", err
	}

	var version string
	sc := bufio.NewScanner(strings.NewReader(string(out)))
	c := 0
	for sc.Scan() {
		if c > 0 {
			break
		}

		if l := sc.Text(); l != "" {
			sp := strings.Split(strings.TrimSpace(sc.Text()), " ")
			full := strings.Split(sp[1], ".")
			version = fmt.Sprintf("%s.%s", full[0], full[1])
		}

		c = c + 1
	}

	return version, nil
}
