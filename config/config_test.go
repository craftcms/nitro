package config

import (
	"testing"

	"github.com/spf13/viper"
)

func TestGetInt(t *testing.T) {
	type args struct {
		key  string
		flag int
	}
	tests := []struct {
		name       string
		keyToSet   string
		valueToSet interface{}
		args       args
		want       int
	}{
		{
			name: "can get the flag when viper is not set",
			args: args{
				key:  "some.key",
				flag: 4,
			},
			want: 4,
		},
		{
			name:       "can get the flag when viper is set",
			keyToSet:   "some.key",
			valueToSet: 5,
			args: args{
				key:  "some.key",
				flag: 0,
			},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.keyToSet != "" {
				viper.Set(tt.keyToSet, tt.valueToSet)
			}

			if got := GetInt(tt.args.key, tt.args.flag); got != tt.want {
				t.Errorf("GetString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetString(t *testing.T) {
	type args struct {
		key  string
		flag string
	}
	tests := []struct {
		name       string
		keyToSet   string
		valueToSet interface{}
		args       args
		want       string
	}{
		{
			name: "can get the flag when viper is not set",
			args: args{
				key:  "some.key",
				flag: "value",
			},
			want: "value",
		},
		{
			name:       "can get the flag when viper is set",
			keyToSet:   "some.key",
			valueToSet: "thevalue",
			args: args{
				key:  "some.key",
				flag: "",
			},
			want: "thevalue",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.keyToSet != "" {
				viper.Set(tt.keyToSet, tt.valueToSet)
			}

			if got := GetString(tt.args.key, tt.args.flag); got != tt.want {
				t.Errorf("GetString() = %v, want %v", got, tt.want)
			}
		})
	}
}
