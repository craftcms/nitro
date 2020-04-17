package testdata

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"

	"github.com/craftcms/nitro/config"
)

func TestCanSaveConfigProperly(t *testing.T) {
	// set the config file
	t.Log("TODO")

	originalFilePath, err := filepath.Abs("configs/full-example.yaml")
	if err != nil {
		t.Fatal(err)
	}
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile(originalFilePath)
	if err := v.ReadInConfig(); err != nil {
		t.Fatal(err)
	}

	// marshal the original file
	var originalConfigMarshal config.Config
	f, err := ioutil.ReadFile(originalFilePath)
	if err != nil {
		t.Fatal(err)
	}
	if err := yaml.Unmarshal(f, &originalConfigMarshal); err != nil {
		t.Fatal(err)
	}

	t.Log(originalConfigMarshal)
	// remove a site
	// remove a mount
	// save to a temp config
	// compare the two files
	// cleanup
}
