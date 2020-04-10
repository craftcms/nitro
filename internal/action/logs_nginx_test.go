package action

import (
	"reflect"
	"testing"
)

func TestLogsNginx(t *testing.T) {
	type args struct {
		name string
		kind string
	}
	tests := []struct {
		name    string
		args    args
		want    *Action
		wantErr bool
	}{
		{
			name: "can show all logs for nginx",
			args: args{
				name: "somename",
				kind: "",
			},
			want: &Action{
				Type:       "exec",
				UseSyscall: false,
				Args:       []string{"exec", "somename", "--", "sudo", "tail", "-f", "/var/log/nginx/access.log", "-f", "/var/log/nginx/error.log"},
			},
			wantErr: false,
		},
		{
			name: "can show all access logs for nginx",
			args: args{
				name: "somename",
				kind: "access",
			},
			want: &Action{
				Type:       "exec",
				UseSyscall: false,
				Args:       []string{"exec", "somename", "--", "sudo", "tail", "-f", "/var/log/nginx/access.log"},
			},
			wantErr: false,
		},
		{
			name: "can show all error logs for nginx",
			args: args{
				name: "somename",
				kind: "error",
			},
			want: &Action{
				Type:       "exec",
				UseSyscall: false,
				Args:       []string{"exec", "somename", "--", "sudo", "tail", "-f", "/var/log/nginx/error.log"},
			},
			wantErr: false,
		},
		{
			name: "missing name returns error",
			args: args{
				name: "",
				kind: "access",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LogsNginx(tt.args.name, tt.args.kind)
			if (err != nil) != tt.wantErr {
				t.Errorf("LogsNginx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LogsNginx() got = %v, want %v", got, tt.want)
			}
		})
	}
}
