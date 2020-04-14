package nitro

import (
	"reflect"
	"testing"
)

func TestLogsDocker(t *testing.T) {
	type args struct {
		name      string
		container string
	}
	tests := []struct {
		name    string
		args    args
		want    *Action
		wantErr bool
	}{
		{
			name: "can show logs for all containers",
			args: args{
				name:      "somename",
				container: "container",
			},
			want: &Action{
				Type:       "exec",
				UseSyscall: false,
				Args:       []string{"exec", "somename", "--", "docker", "logs", "container", "-f"},
			},
			wantErr: false,
		},
		{
			name: "missing name returns error",
			args: args{
				name:      "",
				container: "container",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "missing container returns error",
			args: args{
				name:      "somename",
				container: "",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LogsDocker(tt.args.name, tt.args.container)
			if (err != nil) != tt.wantErr {
				t.Errorf("LogsDocker() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LogsDocker() got = %v, want %v", got, tt.want)
			}
		})
	}
}
