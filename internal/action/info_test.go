package action

import (
	"reflect"
	"testing"
)

func TestInfo(t *testing.T) {
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
			name: "can get a machines info",
			args: args{name: "systemname"},
			want: &Action{
				Type:       "info",
				UseSyscall: false,
				Args:       []string{"info", "systemname"},
			},
			wantErr: false,
		},
		{
			name:    "can get a machines info",
			args:    args{name: ""},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Info(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Info() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Info() got = %v, want %v", got, tt.want)
			}
		})
	}
}
