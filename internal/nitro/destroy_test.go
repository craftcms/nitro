package nitro

import (
	"reflect"
	"testing"
)

func TestDestroy(t *testing.T) {
	type args struct {
		name  string
		force bool
	}
	tests := []struct {
		name    string
		args    args
		want    *Action
		wantErr bool
	}{
		{
			name: "can destroy a machine permanently",
			args: args{
				name:  "ispermanent",
			},
			want: &Action{
				Type:       "delete",
				UseSyscall: false,
				Args:       []string{"delete", "ispermanent", "-p"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Destroy(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Destroy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Destroy() got = %v, want %v", got, tt.want)
			}
		})
	}
}
