package config

type Site struct {
	Hostname string   `yaml:"hostname"`
	Webroot  string   `yaml:"webroot"`
	Aliases  []string `yaml:"aliases,omitempty"`
}
