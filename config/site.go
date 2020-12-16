package config

import (
	"path/filepath"
	"strings"
)

// Site represents a web application. It has a hostname, aliases (which
// are alternate domains), the local path to the site, additional mounts
// to add to the container, and the directory the index.php is located.
type Site struct {
	Hostname string   `mapstructure:"hostname" yaml:"hostname"`
	Aliases  []string `mapstructure:"aliases,omitempty" yaml:"aliases,omitempty"`
	Path     string   `mapstructure:"path" yaml:"path"`
	PHP      string   `mapstructure:"php" yaml:"php"`
	Dir      string   `mapstructure:"dir" yaml:"dir"`
	Xdebug   bool     `mapstructure:"xdebug" yaml:"xdebug"`
}

// GetAbsPath gets the directory for a site.Path,
// It is used to create the mount for a sites
// container.
func (s *Site) GetAbsPath(home string) (string, error) {
	return s.cleanPath(home, s.Path)
}

func (s *Site) cleanPath(home, path string) (string, error) {
	p := path
	if strings.Contains(p, "~") {
		p = strings.Replace(p, "~", home, -1)
	}

	abs, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}

	return filepath.Clean(abs), nil
}
