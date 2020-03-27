package config

import (
	"flag"

	"github.com/spf13/viper"
)

// ReadFile takes a file and default arguments. The override arguments are then passed in as
// from the cli as flags or global defaults such as CPU, Memory, and Disk.
func ReadFile(file string, overrides map[string]interface{}) (*viper.Viper, error) {
	v := viper.New()

	// set all of the overrides
	for key, value := range overrides {
		v.SetDefault(key, value)
	}

	// set the config file name
	v.SetConfigName(file)
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	// try reading the file
	err := v.ReadInConfig()

	return v, err
}

func ReadFileWithFlags(file string, overrides *flag.FlagSet) *viper.Viper {
	v := viper.New()

	if overrides != nil {
		overrides.VisitAll(func(f *flag.Flag) {
			v.SetDefault(f.Name, f.Value)
		})
	}

	v.SetConfigFile(file)
	v.SetConfigType("yaml")
	v.AddConfigPath("$HOME/.nitro/")

	// explicitly ignore the reading of the config. This throws an error when the file is not found
	// and nitro does not require a file on init.
	_ = v.ReadInConfig()

	return v
}
