package config

type Site struct {
	Hostname string   `yaml:"hostname"`
	Webroot  string   `yaml:"webroot"`
	Aliases  []string `yaml:"aliases,omitempty"`
}

func (s *Site) IsExact(site Site) bool {
	if s.Hostname == site.Hostname && s.Webroot == site.Webroot {
		return true
	}

	return false
}
