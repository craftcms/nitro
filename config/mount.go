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

func (m *Mount) IsExact(dest string) bool {
	if m.Dest == dest {
		return true
	}

	return false
}

func (m *Mount) IsParent(dest string) bool {
	if strings.Contains(dest, m.Dest) {
		return true
	}

	return false
}
