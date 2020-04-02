package nitro

import (
	"reflect"
	"testing"
)

func TestConfig_Parse(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name string
		args args
		want *Config
	}{
		{
			name: "load complete nitro",
			args: args{file: "./testdata/nitro-partial.yaml"},
			want: &Config{
				Name:   "defaults",
				CPU:    2,
				Memory: "4G",
				Disk:   "40G",
				Database: struct {
					Engine  string `yaml:"engine"`
					Version string `yaml:"version"`
				}{
					Engine:  "mysql",
					Version: "5.7",
				},
			},
		},
		{
			name: "load complete nitro",
			args: args{file: "./testdata/nitro.yaml"},
			want: &Config{
				Name:   "some-name",
				CPU:    2,
				Memory: "4G",
				Disk:   "33G",
				Database: struct {
					Engine  string `yaml:"engine"`
					Version string `yaml:"version"`
				}{
					Engine:  "mysql",
					Version: "5.7",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{}
			if got := c.Parse(tt.args.file); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\ngot:\n%v, \nwant:\n%v", got, tt.want)
			}
		})
	}
}
