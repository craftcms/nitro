package action

import (
	"reflect"
	"testing"
)

func TestUnmount(t *testing.T) {
	type args struct {
		name string
		site string
	}
	tests := []struct {
		name    string
		args    args
		want    *Action
		wantErr bool
	}{
		{
			name: "valid args return action",
			args: args{name: "somename", site: "example.test"},
			want: &Action{
				Type:       "umount",
				UseSyscall: false,
				Args:       []string{"umount", "somename:/app/sites/example.test"},
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
			got, err := Unmount(tt.args.name, tt.args.site)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unmount() got = %v, want %v", got, tt.want)
			}
		})
	}
}
