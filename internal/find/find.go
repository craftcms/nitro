package find

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
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

// docker container ls --format '{{.Names}}'
func AllDatabases(f Finder) ([]config.Database, error) {
	out, err := f.Output()
	if err != nil {
		return nil, err
	}

	var databases []config.Database
	sc := bufio.NewScanner(strings.NewReader(string(out)))
	for sc.Scan() {
		if strings.Contains(sc.Text(), "mysql") || strings.Contains(sc.Text(), "postgres") {
			sp := strings.Split(sc.Text(), "_")
			db := config.Database{
				Engine:  strings.TrimLeft(sp[0], "'"),
				Version: sp[1],
				Port:    strings.TrimRight(sp[2], "'"),
			}
			databases = append(databases, db)
		}
	}

	return databases, nil
}
