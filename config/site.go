package config


// Site is the representation of a virtual host in nitro.
type Site struct {
	Hostname string   `yaml:"hostname"`
	Webroot  string   `yaml:"webroot"`
	Aliases  []string `yaml:"aliases,omitempty"`
}

// IsExact verifies the current site and the provided
// have the same webroot and hostname. It does not
// currently check aliases
func (s *Site) IsExact(site Site) bool {
	if s.Hostname == site.Hostname && s.Webroot == site.Webroot {
		return true
	}

	return false
}
