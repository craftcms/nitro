package config

import "fmt"

// Database is the struct used to represent a database engine
// that is a combination of a engine (e.g. mariadb, mysql, or
// postgresl), the version number, and the port. The engine
// and version are directly related to the official docker
// images on the docker hub.
type Database struct {
	Engine  string `yaml:"engine,omitempty"`
	Version string `yaml:"version,omitempty"`
	Port    string `yaml:"port,omitempty"`
}

// GetHostname returns a friendly and predictable name for a database
// container. It is used for accessing a database by hostname. For
// example, mysql-8.0-3306 would be the hostname to use in the .env
// for DB_HOST.
func (d *Database) GetHostname() (string, error) {
	if d.Engine == "" || d.Version == "" || d.Port == "" {
		return "", fmt.Errorf("the engine, version, and port must be defined for the database")
	}

	return fmt.Sprintf("%s-%s-%s", d.Engine, d.Version, d.Port), nil
}
