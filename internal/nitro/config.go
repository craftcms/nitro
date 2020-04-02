package nitro

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
)

type Config struct {
	Name     string `yaml:"name"`
	CPU      int    `yaml:"cpu"`
	Memory   string `yaml:"memory"`
	Disk     string `yaml:"disk"`
	Database struct {
		Engine  string `yaml:"engine"`
		Version string `yaml:"version"`
	} `yaml:"database"`
}

func (c *Config) Parse(file string) *Config {
	var conf Config

	// if there is a file load it and unmarshal
	if file != "" {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return nil
		}
		_ = yaml.Unmarshal(data, &conf)
	}

	if conf.CPU == 0 {
		conf.CPU = 2
	}

	if conf.Memory == "" {
		conf.Memory = "4G"
	}

	if conf.Disk == "" {
		conf.Disk = "40G"
	}

	return &conf
}
