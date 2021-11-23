package containerlabels

import "testing"

func TestIsServiceContainer(t *testing.T) {
	type args struct {
		labels map[string]string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "returns true if the labels are for a service",
			args: args{labels: map[string]string{
				Type: "dynamodb",
			}},
			want: true,
		},
		{
			name: "returns true if the labels are for a service",
			args: args{labels: map[string]string{
				Host: "anything",
			}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsServiceContainer(tt.args.labels); got != tt.want {
				t.Errorf("IsServiceContainer() = %v, want %v", got, tt.want)
			}
		})
	}
}
