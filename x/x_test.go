package x

import (
	"reflect"
	"testing"
)

func TestInit(t *testing.T) {
	type args struct {
		machine   string
		cpus      string
		memory    string
		disk      string
		php       string
		dbEngine  string
		dbVersion string
	}
	tests := []struct {
		name string
		args args
		want Command
	}{
		{
			name: "testing the values",
			args: args{
				machine:   "some-machine",
				cpus:      "16",
				memory:    "64G",
				disk:      "160G",
				php:       "7.4",
				dbEngine:  "mysql",
				dbVersion: "5.7",
			},
			want: Command{
				Description: "launch and init",
				Machine:     "some-machine",
				Args: map[string][]string{
					"launch":  {"launch", "--name", "some-machine", "--cpus", "16", "--mem", "64G", "--disk", "160G", "--cloud-init", "-"},
					"install": {"exec", "some-machine", "--", "sudo", "apt", "install", "-y", "php7.4", "php7.4-mbstring", "php7.4-cli", "php7.4-curl", "php7.4-fpm", "php7.4-gd", "php7.4-intl", "php7.4-json", "php7.4-mysql", "php7.4-opcache", "php7.4-pgsql", "php7.4-zip", "php7.4-xml", "php-xdebug", "php-imagick"},
					"docker":  {"exec", "some-machine", "--", "docker", "run", "-d", "--restart=always", "-p", "3306:3306", "-e", "MYSQL_ROOT_PASSWORD=nitro", "-e", "MYSQL_DATABASE=nitro", "-e", "MYSQL_USER=nitro", "-e", "MYSQL_PASSWORD=nitro", "mysql:5.7"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Init(tt.args.machine, tt.args.cpus, tt.args.memory, tt.args.disk, tt.args.php, tt.args.dbEngine, tt.args.dbVersion)

			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("wanted: \n%v \ngot: \n%v", tt.want, got)
			}

			// TODO might delete this, but it was useful
			// loops over all of the args
			//for i, _ := range tt.want.Args {
			//	if !reflect.DeepEqual(tt.want.Args[i], got.Args[i]) {
			//		t.Errorf("wanted: \n%v \ngot: \n%v", tt.want.Args[i], got.Args[i])
			//	}
			//}
		})
	}
}
