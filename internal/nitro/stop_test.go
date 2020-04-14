package nitro

import (
	"reflect"
	"testing"
)

func TestStop(t *testing.T) {
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
			name: "valid args returns nitro",
			args: args{name: "somename"},
			want: &Action{
				Type:       "stop",
				UseSyscall: false,
				Args:       []string{"stop", "somename"},
			},
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
			got, err := Stop(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Stop() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Stop() got = %v, want %v", got, tt.want)
			}
		})
	}
}
