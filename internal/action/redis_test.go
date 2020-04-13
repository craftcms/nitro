package action

import (
	"reflect"
	"testing"
)

func TestRedis(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    *Action
		wantErr bool
	}{
		{
			name: "valid args returns action",
			args: args{name: "somename"},
			want: &Action{
				Type:       "exec",
				UseSyscall: false,
				Args:       []string{"exec", "somename", "--", "redis-cli"},
			},
			wantErr: false,
		},
		{
			name:    "invalid name returns error",
			args:    args{name: ""},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Redis(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Redis() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Redis() got = %v, want %v", got, tt.want)
			}
		})
	}
}
