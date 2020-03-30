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
			want: []string{"exec", "some-name", "--", "docker", "run", "-d", "--restart=always", "postgres" + ":" + "11.5", "-p", "5432:5432", "-e", "POSTGRES_PASSWORD=nitro", "-e", "POSTGRES_USER=nitro", "-e", "POSTGRES_DB=nitro"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DockerRunDatabase(tt.args.name, tt.args.engine, tt.args.version); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DockerRunDatabase() = %v, want %v", got, tt.want)
			}
		})
	}
}
