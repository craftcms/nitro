package config

type Database struct {
	Engine  string `yaml:"engine"`
	Version string `yaml:"version"`
	Port    string `yaml:"port"`
}
