package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestDatabase_GetHostname(t *testing.T) {
	type fields struct {
		Engine    string
		Version   string
		Port      string
		Ephemeral bool
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name:    "can get the hostname for a database container",
			fields:  fields{Engine: "mysql", Version: "5.7", Port: "3306"},
			want:    "mysql-5.7-3306",
			wantErr: false,
		},
		{
			name:    "empty values return an error",
			fields:  fields{Engine: "mysql", Port: "3306"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Database{
				Engine:  tt.fields.Engine,
				Version: tt.fields.Version,
				Port:    tt.fields.Port,
			}
			got, err := d.GetHostname()
			if (err != nil) != tt.wantErr {
				t.Errorf("Database.GetHostname() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Database.GetHostname() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	// get the working dir for the test path
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	testdir := filepath.Join(wd, "testdata", ".nitro")

	type args struct {
		home string
		env  string
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "can load a config file",
			args: args{
				home: "testdata",
				env:  "nitro-test",
			},
			want: &Config{
				file: filepath.Join(testdir, "nitro-test"+".yaml"),
				Blackfire: Blackfire{
					ServerID:    "my-id",
					ServerToken: "my-token",
				},
				Databases: []Database{
					{
						Engine:  "mysql",
						Version: "8.0",
						Port:    "3306",
					},
					{
						Engine:  "postgres",
						Version: "13",
						Port:    "5432",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Load(tt.args.home, tt.args.env)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() = %v, want %v", got, tt.want)
			}
		})
	}
}
