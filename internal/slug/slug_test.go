package slug

import "testing"

func TestGenerate(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "spaces are replaced with _",
			args: args{s: "this is a space"},
			want: "this_is_a_space",
		},
		{
			name: "special chars are replaced with _",
			args: args{s: "this is & special"},
			want: "this_is_special",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Generate(tt.args.s); got != tt.want {
				t.Errorf("Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}
