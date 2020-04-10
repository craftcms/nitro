package config

import "github.com/spf13/viper"

type Config struct {
	Name     string
	PHP      string
	Database []Database
	Sites    []Site
}

func GetString(key, flag string) string {
	if viper.IsSet(key) && flag == "" {
		return viper.GetString(key)
	}

	return flag
}

func GetInt(key string, flag int) int {
	if viper.IsSet(key) && flag == 0 {
		return viper.GetInt(key)
	}

	return flag
}
