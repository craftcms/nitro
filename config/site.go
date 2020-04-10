package config

type Site struct {
	Domain  string `yaml:"domain"`
	Path    string `yaml:"path"`
	Docroot string `yaml:"docroot"`
}
