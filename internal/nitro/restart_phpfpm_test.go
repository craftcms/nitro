package nitro

import (
	"reflect"
	"testing"
)

func TestRestartPhpFpm(t *testing.T) {
	type args struct {
		name string
		php  string
	}
	tests := []struct {
		name    string
		args    args
		want    *Action
		wantErr bool
	}{
		{
			name: "valid args return nitro",
			args: args{name: "somename", php: "7.4"},
			want: &Action{
				Type:       "exec",
				Output:     "Restarting php-fpm 7.4",
				UseSyscall: false,
				Args:       []string{"exec", "somename", "--", "sudo", "service", "php7.4-fpm", "restart"},
			},
			wantErr: false,
		},
		{
			name:    "invalid name returns error",
			args:    args{name: "", php: "7.4"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid php returns error",
			args:    args{name: "somename", php: "7.9"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RestartPhpFpm(tt.args.name, tt.args.php)
			if (err != nil) != tt.wantErr {
				t.Errorf("RestartPhpFpm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RestartPhpFpm() got = %v, want %v", got, tt.want)
			}
		})
	}
}
