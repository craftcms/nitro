package config

import "github.com/spf13/viper"

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
