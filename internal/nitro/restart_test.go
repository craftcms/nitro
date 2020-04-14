package nitro

import (
	"reflect"
	"testing"
)

func TestRestart(t *testing.T) {
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
			name: "valid args return nitro",
			args: args{name: "somename"},
			want: &Action{
				Type:       "restart",
				UseSyscall: false,
				Args:       []string{"restart", "somename"},
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
			got, err := Restart(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Restart() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Restart() got = %v, want %v", got, tt.want)
			}
		})
	}
}
