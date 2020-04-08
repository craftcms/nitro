package action

import (
	"reflect"
	"testing"
)

func TestCreateDatabaseContainer(t *testing.T) {
	type args struct {
		name    string
		engine  string
		version string
		port    string
	}
	tests := []struct {
		name    string
		args    args
		want    *Action
		wantErr bool
	}{
		{
			name: "create mysql 5.7",
			args: args{
				name:    "machinename",
				engine:  "mysql",
				version: "5.7",
				port:    "3306",
			},
			want: &Action{
				Type:       "exec",
				UseSyscall: false,
				Args:       []string{"exec", "machinename", "--", "docker", "run", "-v", "mysql_5.7_3306:/var/lib/mysql", "--name", "mysql_5.7_3306", "-d", "--restart=always", "-p", "3306:3306", "mysql:5.7"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateDatabaseContainer(tt.args.name, tt.args.engine, tt.args.version, tt.args.port)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateDatabaseContainer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateDatabaseContainer() got = \n%v, \nwant \n%v", got, tt.want)
			}
		})
	}
}

func TestCreateDatabaseVolume(t *testing.T) {
	type args struct {
		name    string
		engine  string
		version string
		port    string
	}
	tests := []struct {
		name    string
		args    args
		want    *Action
		wantErr bool
	}{
		{
			name: "create mysql 5.7 volume",
			args: args{
				name:    "somename",
				engine:  "mysql",
				version: "5.7",
				port:    "3306",
			},
			want: &Action{
				Type:       "exec",
				UseSyscall: false,
				Args:       []string{"exec", "somename", "--", "docker", "volume", "create", "mysql_5.7_3306"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateDatabaseVolume(tt.args.name, tt.args.engine, tt.args.version, tt.args.port)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateDatabaseVolume() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateDatabaseVolume() got = \n%v, \nwant \n%v", got, tt.want)
			}
		})
	}
}
