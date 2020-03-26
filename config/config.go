package config

import "github.com/spf13/viper"

func ReadFile(file string, defaults map[string]interface{}) (*viper.Viper, error) {
	v := viper.New()

	// set all of the defaults
	for key, value := range defaults {
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
