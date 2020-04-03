package scripts

import (
	"reflect"
	"testing"
)

func TestRunDatabase(t *testing.T) {
	type args struct {
		name    string
		engine  string
		version string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "can run postgres",
			args: args{
				name:    "some-name",
				engine:  "postgres",
				version: "11.5",
			},
			want: []string{"exec", "some-name", "--", "docker", "run", "-v", "/opt/nitro/volumes/postgres:/var/lib/postgresql/data", "--name", "nitro_postgres_11.5", "-d", "--restart=always", "-p", "5432:5432", "-e", "POSTGRES_PASSWORD=nitro", "-e", "POSTGRES_USER=nitro", "-e", "POSTGRES_DB=nitro", "postgres" + ":" + "11.5"},
		},
		{
			name: "can run mysql",
			args: args{
				name:    "some-name",
				engine:  "mysql",
				version: "5.7",
			},
			want: []string{"exec", "some-name", "--", "docker", "run", "-v", "/opt/nitro/volumes/mysql:/var/lib/mysql", "--name", "nitro_mysql_5.7", "-d", "--restart=always", "-p", "3306:3306", "-e", "MYSQL_ROOT_PASSWORD=nitro", "-e", "MYSQL_DATABASE=nitro", "-e", "MYSQL_USER=nitro", "-e", "MYSQL_PASSWORD=nitro", "mysql" + ":" + "5.7"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DockerRunDatabase(tt.args.name, tt.args.engine, tt.args.version); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DockerRunDatabase() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}
