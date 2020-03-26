package config

import (
	"testing"
)

func TestReadFile(t *testing.T) {
	// Arrange
	file := "./testdata/nitro-full.yaml"
	overrides := map[string]interface{}{
		"cpu":    2,
		"disk":   "20G",
		"memory": "2G",
	}

	// Act
	config, err := ReadFile(file, overrides)
	if err != nil {
		t.Error(err)
	}

	// Assert
	if config.GetString("name") != "nitro-server" {
		t.Errorf("expected %q to be %q; got %q instead", "name", "nitro-server", config.GetString("name"))
	}
	// check if the defaults are set
	if config.GetInt("cpu") != 2 {
		t.Errorf("expected %q to be %q; got %q instead", "cpu", 2, config.GetInt("cpu"))
	}
	if config.GetString("memory") != "2G" {
		t.Errorf("expected %q to be %q; got %q instead", "memory", "2G", config.GetString("memory"))
	}
	if config.GetString("disk") != "20G" {
		t.Errorf("expected %q to be %q; got %q instead", "disk", "20G", config.GetString("disk"))
	}
}
