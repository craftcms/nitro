package action

import (
	"reflect"
	"testing"
)

func TestEnableXdebug(t *testing.T) {
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
			name: "valid args return action",
			args: args{name: "somename", php: "7.4"},
			want: &Action{
				Type:       "exec",
				UseSyscall: false,
				Args:       []string{"exec", "somename", "--", "sudo", "phpenmod", "-v", "7.4", "xdebug"},
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
			got, err := EnableXdebug(tt.args.name, tt.args.php)
			if (err != nil) != tt.wantErr {
				t.Errorf("EnableXdebug() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EnableXdebug() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDisableXdebug(t *testing.T) {
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
			name: "valid args return action",
			args: args{name: "somename", php: "7.4"},
			want: &Action{
				Type:       "exec",
				UseSyscall: false,
				Args:       []string{"exec", "somename", "--", "sudo", "phpdismod", "-v", "7.4", "xdebug"},
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
			got, err := DisableXdebug(tt.args.name, tt.args.php)
			if (err != nil) != tt.wantErr {
				t.Errorf("DisableXdebug() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DisableXdebug() got = %v, want %v", got, tt.want)
			}
		})
	}
}
