package suggest

import "testing"

func TestNumberOfCPUs(t *testing.T) {
	type args struct {
		num int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "all others returns four",
			args: args{num: 16},
			want: "4",
		},
		{
			name: "eight cpus returns two",
			args: args{num: 8},
			want: "2",
		},
		{
			name: "six cpus returns four",
			args: args{num: 6},
			want: "4",
		},
		{
			name: "four cpus returns two",
			args: args{num: 4},
			want: "1",
		},
		{
			name: "two cpus returns one",
			args: args{num: 2},
			want: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NumberOfCPUs(tt.args.num); got != tt.want {
				t.Errorf("NumberOfCPUs() = %v, want %v", got, tt.want)
			}
		})
	}
}
