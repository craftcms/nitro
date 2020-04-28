package config

import "fmt"

type Database struct {
	Engine  string `yaml:"engine"`
	Version string `yaml:"version"`
	Port    string `yaml:"port"`
}

// Name converts a database into a name used for the container
func (d *Database) Name() string {
	return fmt.Sprintf("%s_%s_%s", d.Engine, d.Version, d.Port)
}
