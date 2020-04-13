package config

type Mount struct {
	Source      string `yaml:"source"`
	Destination string `yaml:"destination"`
}
