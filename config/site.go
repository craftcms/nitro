package config

import (
	"path/filepath"
	"strings"
)

// Site represents a web application. It has a hostname, aliases (which
// are alternate domains), the local path to the site, additional mounts
// to add to the container, and the directory the index.php is located.
type Site struct {
	Hostname string   `yaml:"hostname,omitempty"`
	Aliases  []string `yaml:"aliases,omitempty"`
	Path     string   `yaml:"path,omitempty"`
	Mounts   []string `yaml:"mounts,omitempty"`
	PHP      string   `yaml:"php,omitempty"`
	Dir      string   `yaml:"dir,omitempty"`
}

// GetAbsPath gets the directory for a site.Path,
// It is used to create the mount for a sites
// container.
func (s *Site) GetAbsPath(home string) (string, error) {
	return s.cleanPath(home, s.Path)
}

// GetAbsMountPaths gets the directory for a site.Mounts,
// It is used to create the additional mounts for a sites
// container.
func (s *Site) GetAbsMountPaths(home string) (map[string]string, error) {
	mnts := make(map[string]string)

	for _, m := range s.Mounts {
		// split the
		sp := strings.Split(m, ":")
		src := sp[0]
		dest := sp[1]

		p, err := s.cleanPath(home, src)
		if err != nil {
			return mnts, err
		}

		mnts[p] = dest
	}

	return mnts, nil
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
