package containerargs

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		args    args
		want    *Command
		wantErr bool
	}{
		{
			name: "passing -- returns the container name and args",
			args: args{args: []string{"sitename.nitro", "--", "install"}},
			want: &Command{
				Container: "sitename.nitro",
				Args:      []string{"install"},
			},
			wantErr: false,
		},
		{
			name: "returns the args if no delimiter was provided",
			args: args{args: []string{"install"}},
			want: &Command{
				Args: []string{"install"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
