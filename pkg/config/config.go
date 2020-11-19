package config

import "fmt"

type Config struct {
	Extensions []string   `yaml:"exts,omitempty"`
	Blackfire  Blackfire  `yaml:"blackfire,omitempty"`
	Sites      []Site     `yaml:"sites,omitempty"`
	Databases  []Database `yaml:"databases,omitempty"`
}

type Blackfire struct {
	ClientID     string `yaml:"client_id,omitempty"`
	ClientSecret string `yaml:"client_secret,omitempty"`
}

type Site struct {
	Hostname string   `yaml:"hostname,omitempty"`
	Aliases  []string `yaml:"aliases,omitempty"`
	Path     string   `yaml:"path,omitempty"`
	PHP      string   `yaml:"php,omitempty"`
}

type Database struct {
	Engine  string `yaml:"engine,omitempty"`
	Version string `yaml:"version,omitempty"`
	Port    string `yaml:"port,omitempty"`
}

func (d *Database) GetHostname() (string, error) {
	if d.Engine == "" || d.Version == "" || d.Port == "" {
		return "", fmt.Errorf("the engine, version, and port must be defined for the database")
	}

	return fmt.Sprintf("%s-%s-%s", d.Engine, d.Version, d.Port), nil
}
