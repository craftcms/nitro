package config

import (
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func getConfigFile(t *testing.T, file string) Config {
	fp, err := filepath.Abs(file)
	if err != nil {
		t.Fatal(err)
	}
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile(fp)
	if err := v.ReadInConfig(); err != nil {
		t.Fatal(err)
	}

	var cfg Config
	f, err := ioutil.ReadFile(fp)
	if err != nil {
		t.Fatal(err)
	}
	if err := yaml.Unmarshal(f, &cfg); err != nil {
		t.Fatal(err)
	}

	return cfg
}

func TestCanSaveConfigProperly(t *testing.T) {
	// Arrange
	// set the config file
	originalCfgFile := getConfigFile(t, "testdata/configs/full-example.yaml")

	// Act
	// remove a site
	if err := originalCfgFile.RemoveSite("demo-site"); err != nil {
		t.Error(err)
	}
	// remove a mount
	if err := originalCfgFile.RemoveMountBySiteWebroot("/nitro/sites/demo-site/web"); err != nil {
		t.Error(err)
	}
	// save to a temp config
	if err := originalCfgFile.Save("testdata/configs/test-example.yaml"); err != nil {
		t.Error(err)
	}

	// Assert
	// compare the original and saved files
	savedCfgFile := getConfigFile(t, "testdata/configs/test-example.yaml")
	if !reflect.DeepEqual(originalCfgFile, savedCfgFile) {
		t.Errorf("expected configs to be the same, got \n%v \nwant: \n%v", originalCfgFile, savedCfgFile)
	}
	// double check the golden file
	goldenCfgFile := getConfigFile(t, "testdata/configs/golden-full.yaml")
	if !reflect.DeepEqual(goldenCfgFile, savedCfgFile) {
		t.Errorf("expected configs to be the same, got \n%v \nwant: \n%v", goldenCfgFile, savedCfgFile)
	}
}
