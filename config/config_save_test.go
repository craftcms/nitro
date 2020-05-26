package config

import (
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func getConfigFile(t *testing.T, file string) *Config {
	fp, err := filepath.Abs(file)
	if err != nil {
		t.Fatal(err)
	}
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile(fp)
	if err := v.ReadInConfig(); err != nil {
		t.Fatal("error reading", file, "err:", err)
	}

	var cfg Config
	f, err := ioutil.ReadFile(fp)
	if err != nil {
		t.Fatal(err)
	}
	if err := yaml.Unmarshal(f, &cfg); err != nil {
		t.Fatal(err)
	}

	return &cfg
}

func TestCanSaveConfigProperly(t *testing.T) {
	// Arrange
	// set the config file
	originalCfgFile := getConfigFile(t, "testdata/configs/full-example.yaml")

	// Act
	sites := originalCfgFile.GetSites()
	site := sites[0]
	_ = originalCfgFile.FindMountBySiteWebroot(site.Webroot)

	// remove a site
	if err := originalCfgFile.RemoveSite(site.Hostname); err != nil {
		t.Error(err)
	}
	// remove a mount
	if err := originalCfgFile.RemoveMountBySiteWebroot(site.Webroot); err != nil {
		t.Error(err)
	}
	// save to a temp config
	if err := originalCfgFile.Save("testdata/configs/test-example.yaml"); err != nil {
		t.Error(err)
	}

	// Assert
	// compare the original and saved files
	modifiedCfgFile := getConfigFile(t, "testdata/configs/test-example.yaml")
	// double check the golden file
	goldenCfgFile := getConfigFile(t, "testdata/configs/golden-full.yaml")
	// make sure the mounts are the same
	if !reflect.DeepEqual(goldenCfgFile.Mounts, modifiedCfgFile.Mounts) {
		t.Errorf("expected configs mounts to be the same, got \n%v \nwant: \n%v", modifiedCfgFile.Mounts, goldenCfgFile.Mounts)
	}
	t.Log("this test needs to be updated")
	// make sure the sites are the same
	//if !reflect.DeepEqual(goldenCfgFile.Sites, modifiedCfgFile.Sites) {
	//	t.Errorf("expected configs sites to be the same, got \n%p \nwant: \n%p", &modifiedCfgFile.Sites, &goldenCfgFile.Sites)
	//}
	//// make sure the entire config is the same
	//if !reflect.DeepEqual(goldenCfgFile, modifiedCfgFile) {
	//	t.Errorf("expected configs to be the same, got \n%v \nwant: \n%v", goldenCfgFile, modifiedCfgFile)
	//}
}
