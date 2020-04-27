package config

import (
	"reflect"
	"testing"

	"github.com/spf13/viper"
)

func TestGetInt(t *testing.T) {
	type args struct {
		key  string
		flag int
	}
	tests := []struct {
		name       string
		keyToSet   string
		valueToSet interface{}
		args       args
		want       int
	}{
		{
			name: "can get the flag when viper is not set",
			args: args{
				key:  "some.key",
				flag: 4,
			},
			want: 4,
		},
		{
			name:       "can get the flag when viper is set",
			keyToSet:   "some.key",
			valueToSet: 5,
			args: args{
				key:  "some.key",
				flag: 0,
			},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.keyToSet != "" {
				viper.Set(tt.keyToSet, tt.valueToSet)
			}

			if got := GetInt(tt.args.key, tt.args.flag); got != tt.want {
				t.Errorf("GetString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetString(t *testing.T) {
	type args struct {
		key  string
		flag string
	}
	tests := []struct {
		name       string
		keyToSet   string
		valueToSet interface{}
		args       args
		want       string
	}{
		{
			name: "can get the flag when viper is not set",
			args: args{
				key:  "some.key",
				flag: "value",
			},
			want: "value",
		},
		{
			name:       "can get the flag when viper is set",
			keyToSet:   "some.key",
			valueToSet: "thevalue",
			args: args{
				key:  "some.key",
				flag: "",
			},
			want: "thevalue",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.keyToSet != "" {
				viper.Set(tt.keyToSet, tt.valueToSet)
			}

			if got := GetString(tt.args.key, tt.args.flag); got != tt.want {
				t.Errorf("GetString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_RemoveSite(t *testing.T) {
	type fields struct {
		Name      string
		PHP       string
		CPUs      string
		Disk      string
		Memory    string
		Databases []Database
		Sites     []Site
	}
	type args struct {
		site string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Site
		wantErr bool
	}{
		{
			name: "remove a site by its hostname",
			args: args{
				site: "anotherexample.test",
			},
			fields: fields{
				Sites: []Site{
					{
						Hostname: "example.test",
						Webroot:  "web",
					},
					{
						Hostname: "anotherexample.test",
						Webroot:  "web",
					},
					{
						Hostname: "finalexample.test",
						Webroot:  "web",
					},
				},
			},
			want: []Site{
				{
					Hostname: "example.test",
					Webroot:  "web",
				},
				{
					Hostname: "finalexample.test",
					Webroot:  "web",
				},
			},
			wantErr: false,
		},
		{
			name: "sites not in the slice return an error",
			args: args{
				site: "doesnotexist.test",
			},
			fields: fields{
				Sites: []Site{
					{
						Hostname: "example.test",
						Webroot:  "web",
					},
					{
						Hostname: "anotherexample.test",
						Webroot:  "web",
					},
					{
						Hostname: "finalexample.test",
						Webroot:  "web",
					},
				},
			},
			want: []Site{
				{
					Hostname: "example.test",
					Webroot:  "web",
				},
				{
					Hostname: "anotherexample.test",
					Webroot:  "web",
				},
				{
					Hostname: "finalexample.test",
					Webroot:  "web",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				PHP:       tt.fields.PHP,
				CPUs:      tt.fields.CPUs,
				Disk:      tt.fields.Disk,
				Memory:    tt.fields.Memory,
				Databases: tt.fields.Databases,
				Sites:     tt.fields.Sites,
			}

			err := c.RemoveSite(tt.args.site)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveSite() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.want != nil {
				if !reflect.DeepEqual(c.Sites, tt.want) {
					t.Errorf("RemoveSite() got = \n%v, \nwant \n%v", c.Sites, tt.want)
				}
			}
		})
	}
}

func TestConfig_AddMount(t *testing.T) {
	// since paths are different across systems, we automate the ~ directory
	t.Skip("skipping for now, need to update for relative paths")

	type fields struct {
		Name      string
		PHP       string
		CPUs      string
		Disk      string
		Memory    string
		Mounts    []Mount
		Databases []Database
		Sites     []Site
	}
	type args struct {
		m Mount
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Mount
		wantErr bool
	}{
		{
			name: "adds a new mount and sets a default dest path for period references",
			args: args{
				m: Mount{
					Source: "./testdata/test-mount",
				},
			},
			want: []Mount{
				{
					Source: "~/go/src/github.com/craftcms/nitro/config/testdata/test-mount",
					Dest:   "/nitro/sites/test-mount",
				},
			},
			wantErr: false,
		},
		{
			name: "adds a new mount and sets a default dest path for non-relative references",
			args: args{
				m: Mount{
					Source: "testdata/test-mount",
				},
			},
			want: []Mount{
				{
					Source: "~/go/src/github.com/craftcms/nitro/config/testdata/test-mount",
					Dest:   "/nitro/sites/test-mount",
				},
			},
			wantErr: false,
		},
		{
			name: "adds a new mount and sets a default dest path for relative",
			args: args{
				m: Mount{
					Source: "~/go/src/github.com/craftcms/nitro/config/testdata/test-mount",
				},
			},
			want: []Mount{
				{
					Source: "~/go/src/github.com/craftcms/nitro/config/testdata/test-mount",
					Dest:   "/nitro/sites/test-mount",
				},
			},
			wantErr: false,
		},
		{
			name: "adds a new mount",
			args: args{
				m: Mount{
					Source: "~/go/src/github.com/craftcms/nitro/config/testdata/test-mount",
					Dest:   "/home/ubuntu/sites",
				},
			},
			want: []Mount{
				{
					Source: "~/go/src/github.com/craftcms/nitro/config/testdata/test-mount",
					Dest:   "/home/ubuntu/sites",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				PHP:       tt.fields.PHP,
				CPUs:      tt.fields.CPUs,
				Disk:      tt.fields.Disk,
				Memory:    tt.fields.Memory,
				Mounts:    tt.fields.Mounts,
				Databases: tt.fields.Databases,
				Sites:     tt.fields.Sites,
			}
			if err := c.AddMount(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("AddMount() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.want != nil {
				if !reflect.DeepEqual(c.Mounts, tt.want) {
					t.Errorf("AddMount() got = \n%v, \nwant \n%v", c.Mounts, tt.want)
				}
			}
		})
	}
}

func TestConfig_AddSite(t *testing.T) {
	type fields struct {
		Name      string
		PHP       string
		CPUs      string
		Disk      string
		Memory    string
		Mounts    []Mount
		Databases []Database
		Sites     []Site
	}
	type args struct {
		site Site
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Site
		wantErr bool
	}{
		{
			name: "adds a new site",
			args: args{
				site: Site{
					Hostname: "craftdev",
				},
			},
			want: []Site{
				{
					Hostname: "craftdev",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				PHP:       tt.fields.PHP,
				CPUs:      tt.fields.CPUs,
				Disk:      tt.fields.Disk,
				Memory:    tt.fields.Memory,
				Mounts:    tt.fields.Mounts,
				Databases: tt.fields.Databases,
				Sites:     tt.fields.Sites,
			}
			if err := c.AddSite(tt.args.site); (err != nil) != tt.wantErr {
				t.Errorf("AddSite() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.want != nil {
				if !reflect.DeepEqual(c.Sites, tt.want) {
					t.Errorf("AddSite() got = \n%v, \nwant \n%v", c.Sites, tt.want)
				}
			}
		})
	}
}

func TestConfig_RemoveMountBySiteWebroot(t *testing.T) {
	type fields struct {
		Name      string
		PHP       string
		CPUs      string
		Disk      string
		Memory    string
		Mounts    []Mount
		Databases []Database
		Sites     []Site
	}
	type args struct {
		webroot string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Mount
		wantErr bool
	}{
		{
			name: "can remove a mount by a site webroot",
			fields: fields{
				Name:   "somename",
				PHP:    "7.4",
				CPUs:   "3",
				Disk:   "20G",
				Memory: "4G",
				Mounts: []Mount{
					{
						Source: "./testdata/test-mount",
						Dest:   "/nitro/sites/testmount",
					},
					{
						Source: "./testdata/test-mount/remove",
						Dest:   "/nitro/sites/remove",
					},
				},
				Sites: []Site{
					{
						Hostname: "keep.test",
						Webroot:  "/nitro/sites/keep/web",
					},
					{
						Hostname: "remove.test",
						Webroot:  "/nitro/sites/testmount/remove/web",
					},
				},
			},
			args: args{webroot: "/nitro/sites/remove/web"},
			want: []Mount{
				{
					Source: "./testdata/test-mount",
					Dest:   "/nitro/sites/testmount",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				PHP:       tt.fields.PHP,
				CPUs:      tt.fields.CPUs,
				Disk:      tt.fields.Disk,
				Memory:    tt.fields.Memory,
				Mounts:    tt.fields.Mounts,
				Databases: tt.fields.Databases,
				Sites:     tt.fields.Sites,
			}
			if err := c.RemoveMountBySiteWebroot(tt.args.webroot); (err != nil) != tt.wantErr {
				t.Errorf("RemoveMountBySiteWebroot() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.want != nil {
				if !reflect.DeepEqual(c.Mounts, tt.want) {
					t.Errorf("RemoveMountBySiteWebroot() got = \n%v, \nwant \n%v", c.Mounts, tt.want)
				}
			}
		})
	}
}

func TestConfig_RemoveSite1(t *testing.T) {
	type fields struct {
		Name      string
		PHP       string
		CPUs      string
		Disk      string
		Memory    string
		Mounts    []Mount
		Databases []Database
		Sites     []Site
	}
	type args struct {
		hostname string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Site
		wantErr bool
	}{
		{
			name: "can remove a site by its hostname",
			fields: fields{
				Name:   "somename",
				PHP:    "7.4",
				CPUs:   "3",
				Disk:   "20G",
				Memory: "4G",
				Sites: []Site{
					{
						Hostname: "keep.test",
						Webroot:  "/nitro/sites/keep/web",
					},
					{
						Hostname: "remove.test",
						Webroot:  "/nitro/sites/remove/web",
					},
				},
			},
			args: args{hostname: "remove.test"},
			want: []Site{
				{
					Hostname: "keep.test",
					Webroot:  "/nitro/sites/keep/web",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				PHP:       tt.fields.PHP,
				CPUs:      tt.fields.CPUs,
				Disk:      tt.fields.Disk,
				Memory:    tt.fields.Memory,
				Mounts:    tt.fields.Mounts,
				Databases: tt.fields.Databases,
				Sites:     tt.fields.Sites,
			}
			if err := c.RemoveSite(tt.args.hostname); (err != nil) != tt.wantErr {
				t.Errorf("RemoveSite() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.want != nil {
				if !reflect.DeepEqual(c.Sites, tt.want) {
					t.Errorf("RemoveSite() got = \n%v, \nwant \n%v", c.Sites, tt.want)
				}
			}
		})
	}
}

func TestConfig_RenameSite(t *testing.T) {
	type fields struct {
		Mounts    []Mount
		Databases []Database
		Sites     []Site
	}
	type args struct {
		site     Site
		hostname string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Site
		wantErr bool
	}{
		{
			name: "remove a site my hostname",
			args: args{
				site: Site{
					Hostname: "old.test",
					Webroot:  "/nitro/sites/old.test",
				},
				hostname: "new.test",
			},
			fields: fields{
				Sites: []Site{
					{
						Hostname: "old.test",
						Webroot:  "/nitro/sites/old.test",
					},
					{
						Hostname: "keep.test",
						Webroot:  "/nitro/sites/keep.test",
					},
				},
			},
			want: []Site{
				{
					Hostname: "new.test",
					Webroot:  "/nitro/sites/new.test",
				},
				{
					Hostname: "keep.test",
					Webroot:  "/nitro/sites/keep.test",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Mounts:    tt.fields.Mounts,
				Databases: tt.fields.Databases,
				Sites:     tt.fields.Sites,
			}
			if err := c.RenameSite(tt.args.site, tt.args.hostname); (err != nil) != tt.wantErr {
				t.Errorf("RenameSite() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.want != nil {
				if !reflect.DeepEqual(c.Sites, tt.want) {
					t.Errorf("RenameSite() got sites = \n%v, \nwant \n%v", c.Sites, tt.want)
				}
			}
		})
	}
}

func TestConfig_MountExists(t *testing.T) {
	type fields struct {
		PHP       string
		CPUs      string
		Disk      string
		Memory    string
		Mounts    []Mount
		Databases []Database
		Sites     []Site
	}
	type args struct {
		dest string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "existing mounts return true",
			fields: fields{
				Mounts: []Mount{
					{
						Source: "./testdata/existing-mount",
						Dest:   "/nitro/sites/example-site",
					},
				},
			},
			args: args{dest: "/nitro/sites/example-site"},
			want: true,
		},
		{
			name: "non-existing mounts return false",
			fields: fields{
				Mounts: []Mount{
					{
						Source: "./testdata/existing-mount",
						Dest:   "/nitro/sites/example-site",
					},
				},
			},
			args: args{dest: "/nitro/sites/nonexistent-site"},
			want: false,
		},
		{
			name: "parent mounts return true",
			fields: fields{
				Mounts: []Mount{
					{
						Source: "./testdata/test-mount",
						Dest:   "/nitro/sites",
					},
					{
						Source: "./testdata/existing-mount",
						Dest:   "/nitro/sites",
					},
				},
			},
			args: args{dest: "/nitro/sites/nonexistent-site"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				PHP:       tt.fields.PHP,
				CPUs:      tt.fields.CPUs,
				Disk:      tt.fields.Disk,
				Memory:    tt.fields.Memory,
				Mounts:    tt.fields.Mounts,
				Databases: tt.fields.Databases,
				Sites:     tt.fields.Sites,
			}
			if got := c.MountExists(tt.args.dest); got != tt.want {
				t.Errorf("MountExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_SiteExists(t *testing.T) {
	type fields struct {
		PHP       string
		CPUs      string
		Disk      string
		Memory    string
		Mounts    []Mount
		Databases []Database
		Sites     []Site
	}
	type args struct {
		site Site
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "exact sites return true",
			fields: fields{
				Sites: []Site{
					{
						Hostname: "iexist.test",
						Webroot:  "/nitro/sites/iexist.test",
					},
				},
			},
			args: args{site: Site{
				Hostname: "iexist.test",
				Webroot:  "/nitro/sites/iexist.test",
			}},
			want: true,
		},
		{
			name: "exact sites return false",
			fields: fields{
				Sites: []Site{
					{
						Hostname: "iexist.test",
						Webroot:  "/nitro/sites/iexist.test",
					},
				},
			},
			args: args{site: Site{
				Hostname: "idontexist.test",
				Webroot:  "/nitro/sites/idontexist.test",
			}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				PHP:       tt.fields.PHP,
				CPUs:      tt.fields.CPUs,
				Disk:      tt.fields.Disk,
				Memory:    tt.fields.Memory,
				Mounts:    tt.fields.Mounts,
				Databases: tt.fields.Databases,
				Sites:     tt.fields.Sites,
			}
			if got := c.SiteExists(tt.args.site); got != tt.want {
				t.Errorf("SiteExists() = %v, want %v", got, tt.want)
			}
		})
	}
}
