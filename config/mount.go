package config

import (
	"strings"

	"github.com/mitchellh/go-homedir"
)

type Mount struct {
	Source string `yaml:"source"`
	Dest   string `yaml:"dest"`
}

func (m *Mount) AbsSourcePath() string {
	home, _ := homedir.Dir()
	return strings.Replace(m.Source, "~", home, 1)
}

func (m *Mount) Exists(webroot string) bool {
	split := strings.Split(webroot, "/")
	path := split[:len(split)-1]
	dest := strings.Join(path, "/")

	if strings.Contains(m.Dest, dest) {
		return true
	}

	return false
}
