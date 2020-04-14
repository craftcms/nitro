package nitro

import (
	"reflect"
	"testing"
)

func TestSSH(t *testing.T) {
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
			args: args{
				name: "somename",
			},
			want: &Action{
				Type:       "shell",
				UseSyscall: true,
				Args:       []string{"shell", "somename"},
			},
			wantErr: false,
		},
		{
			name: "bad name returns error",
			args: args{
				name: "",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SSH(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("SSH() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SSH() got = %v, want %v", got, tt.want)
			}
		})
	}
}
