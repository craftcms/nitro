package command

import (
	"reflect"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/spf13/viper"
)

func testSSHCommand(t testing.TB, runner *SpyRunner, v *viper.Viper) (*cli.MockUi, *SSHCommand) {
	t.Helper()

	if v != nil {
		_ = v.ReadInConfig()
	}

	ui := cli.NewMockUi()
	coreCmd := &CoreCommand{
		UI:     ui,
		Runner: runner,
		Config: v,
	}

	return ui, &SSHCommand{CoreCommand: coreCmd}
}

func TestSSHCommand_Run(t *testing.T) {
	configWithFile := viper.New()
	configWithFile.SetConfigFile("./testdata/nitro.yaml")
	configWithFile.SetConfigType("yaml")
	config := viper.New()
	config.Set("name", "")
	spyRunner := &SpyRunner{}

	tests := []struct {
		name     string
		args     []string
		expected []string
		config   *viper.Viper
		want     int
	}{
		//{
		//	name:     "ssh command gets the right arguments using a config file",
		//	args:     nil,
		//	expected: []string{"shell", "from-config-file"},
		//	config:   configWithFile,
		//	want:     0,
		//},
		{
			name:     "ssh command gets the right arguments using a flag",
			args:     []string{"-name", "some-machine"},
			expected: []string{"shell", "some-machine"},
			config:   config,
			want:     0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, c := testSSHCommand(t, spyRunner, tt.config)

			if got := c.Run(tt.args); got != tt.want {
				t.Errorf("Run() = %v, want %v", got, tt.want)
			}

			if tt.expected != nil {
				if !reflect.DeepEqual(tt.expected, spyRunner.calls) {
					t.Errorf("wanted: \n%v \ngot: \n%v", tt.expected, spyRunner.calls)
				}
			}
		})
	}
}
