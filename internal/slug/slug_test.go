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
			name: "removes spaces from strings",
			args: args{s: "this database"},
			want: "this_database",
		},
		{
			name: "removes spaces and dashses from strings",
			args: args{s: "this database-test"},
			want: "this_database_test",
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
