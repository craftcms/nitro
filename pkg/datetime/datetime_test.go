package datetime

import (
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	dt, err := time.Parse("2006/01/02 15:04:05", "2020/01/02 08:01:02")
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "can view the time in the correct format",
			args: args{t: dt},
			want: "2020-01-02-080102",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Parse(tt.args.t); got != tt.want {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
