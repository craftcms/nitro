package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoad(t *testing.T) {
	// get the working dir for the test path
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	testdir := filepath.Join(wd, "testdata")

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
				home: testdir,
				env:  "nitro-test",
			},
			want: &Config{
				File: filepath.Join(testdir, ".nitro", "nitro-test"+".yaml"),
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
		{
			name: "no env returns error",
			args: args{
				home: testdir,
			},
			wantErr: true,
		},
		{
			name: "missing file returns an error",
			args: args{
				home: testdir,
				env:  "not-here",
			},
			wantErr: true,
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
				// check blackfire
				if !reflect.DeepEqual(got.Blackfire, tt.want.Blackfire) {
					t.Errorf("Load() = \ngot:\n%v,\nwant\n%v", got.Blackfire, tt.want.Blackfire)
				}

				// check databases
				if !reflect.DeepEqual(got.Databases, tt.want.Databases) {
					t.Errorf("Load() = \ngot:\n%v,\nwant\n%v", got.Databases, tt.want.Databases)
				}

				// check services
				if !reflect.DeepEqual(got.Services, tt.want.Services) {
					t.Errorf("Load() = \ngot:\n%v,\nwant\n%v", got.Services, tt.want.Services)
				}

				// check sites
				if !reflect.DeepEqual(got.Sites, tt.want.Sites) {
					t.Errorf("Load() = \ngot:\n%v,\nwant\n%v", got.Sites, tt.want.Sites)
				}

				t.Errorf("Load() = \ngot\n%v,\nwant\n%v", got, tt.want)
			}
		})
	}
}

func TestConfig_EnableXdebug(t *testing.T) {
	type fields struct {
		Blackfire Blackfire
		Databases []Database
		Services  Services
		Sites     []Site
		File      string
	}
	type args struct {
		site string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "can enable xdebug for a site",
			fields: fields{
				Sites: []Site{
					{
						Hostname: "somesite",
						Xdebug:   false,
					},
					{
						Hostname: "anothersite",
						Xdebug:   true,
					},
				},
			},
			args:    args{site: "somesite"},
			wantErr: false,
		},
		{
			name:    "sites that don't exist return an error",
			args:    args{site: "idontexist"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Blackfire: tt.fields.Blackfire,
				Databases: tt.fields.Databases,
				Services:  tt.fields.Services,
				Sites:     tt.fields.Sites,
				File:      tt.fields.File,
			}
			if err := c.EnableXdebug(tt.args.site); (err != nil) != tt.wantErr {
				t.Errorf("Config.EnableXdebug() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_DisableXdebug(t *testing.T) {
	type fields struct {
		Blackfire Blackfire
		Databases []Database
		Services  Services
		Sites     []Site
		File      string
	}
	type args struct {
		site string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "can enable xdebug for a site",
			fields: fields{
				Sites: []Site{
					{
						Hostname: "somesite",
						Xdebug:   false,
					},
					{
						Hostname: "anothersite",
						Xdebug:   true,
					},
				},
			},
			args:    args{site: "somesite"},
			wantErr: false,
		},
		{
			name:    "sites that don't exist return an error",
			args:    args{site: "idontexist"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Blackfire: tt.fields.Blackfire,
				Databases: tt.fields.Databases,
				Services:  tt.fields.Services,
				Sites:     tt.fields.Sites,
				File:      tt.fields.File,
			}
			if err := c.DisableXdebug(tt.args.site); (err != nil) != tt.wantErr {
				t.Errorf("Config.DisableXdebug() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_AddSite(t *testing.T) {
	type fields struct {
		Blackfire Blackfire
		Databases []Database
		Services  Services
		Sites     []Site
		File      string
	}
	type args struct {
		s Site
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "existing path returns an error",
			fields: fields{
				Sites: []Site{
					{
						Hostname: "example",
						Path:     "some/path",
					},
				},
			},
			args: args{
				s: Site{
					Hostname: "nomatch",
					Path:     "some/path",
				},
			},
			wantErr: true,
		},
		{
			name: "existing hostnames returns an error",
			fields: fields{
				Sites: []Site{
					{
						Hostname: "example",
					},
				},
			},
			args: args{
				s: Site{Hostname: "example"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Blackfire: tt.fields.Blackfire,
				Databases: tt.fields.Databases,
				Services:  tt.fields.Services,
				Sites:     tt.fields.Sites,
				File:      tt.fields.File,
			}
			if err := c.AddSite(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("Config.AddSite() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
