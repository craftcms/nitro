package action

import (
	"reflect"
	"testing"
)

func TestLaunch(t *testing.T) {
	type args struct {
		name   string
		cpus   int
		memory string
		disk   string
		input  string
	}
	tests := []struct {
		name    string
		args    args
		want    *Action
		wantErr bool
	}{
		{
			name: "can launch a new instance",
			args: args{
				name:   "machine",
				cpus:   4,
				memory: "2G",
				disk:   "20G",
				input:  "someinput",
			},
			want: &Action{
				Type:       "launch",
				UseSyscall: false,
				Input:      "someinput",
				Args:       []string{"launch", "--name", "machine", "--cpus", "4", "--mem", "2G", "--disk", "20G", "--cloud-init", "-"},
			},
			wantErr: false,
		},
		{
			name: "missing disk param returns an error",
			args: args{
				name:   "machine",
				cpus:   4,
				memory: "2G",
				input:  "someinput",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "missing input param returns an error",
			args: args{
				name:   "machine",
				cpus:   4,
				memory: "2G",
				disk:   "20G",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "missing cpus param returns an error",
			args: args{
				name:   "machine",
				memory: "2G",
				disk:   "20G",
				input:  "someinput",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "missing memory param returns an error",
			args: args{
				name:  "machine",
				cpus:  4,
				disk:  "20G",
				input: "someinput",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "missing name param returns an error",
			args: args{
				cpus:   4,
				memory: "2G",
				disk:   "20G",
				input:  "someinput",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Launch(tt.args.name, tt.args.cpus, tt.args.memory, tt.args.disk, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Launch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Launch() got = \n%v, \nwant \n%v", got, tt.want)
			}
		})
	}
}
