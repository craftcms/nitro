package command

import (
	"reflect"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/spf13/viper"
)

func testInitCommand(t testing.TB, runner *SpyRunner, v *viper.Viper) (*cli.MockUi, *InitCommand) {
	t.Helper()

	ui := cli.NewMockUi()
	coreCmd := &CoreCommand{
		UI:     ui,
		Runner: runner,
		Config: v,
	}

	return ui, &InitCommand{CoreCommand: coreCmd}
}

func TestInitCommand_Synopsis(t *testing.T) {
	// Arrange
	v := viper.New()
	_, c := testInitCommand(t, nil, v)
	expected := "create new machine"

	// Act
	actual := c.Synopsis()

	// Assert
	if actual != expected {
		t.Errorf("expected %q; got %q instead", expected, actual)
	}
}

func TestInitCommand_Help(t *testing.T) {
	// Arrange
	v := viper.New()
	_, c := testInitCommand(t, nil, v)
	expected := strings.TrimSpace(`
Usage: nitro init [options]
  This command starts a nitro virtual machine and will provision the machine with the requested specifications.
  
  Create a new virtual machine and override the default system specifications:
      $ nitro init -name=diesel -cpu=4 -memory=4G -disk=40GB
  
  Create a new virtual machine and with the defaults and skip bootstrapping the machine with the default installations:
      $ nitro init -name=diesel -skip-install
`)

	// Act
	actual := c.Help()

	// Assert
	if actual != expected {
		t.Errorf("expected %q; got %q instead", expected, actual)
	}
}

func TestInitCommand_Flags(t *testing.T) {
	// Arrange
	v := viper.New()
	_, c := testInitCommand(t, nil, v)
	args := []string{"-name=nitro-dev", "-cpu=2", "-memory=2GB", "-disk=20GB", "-skip-install"}

	// Act
	if err := c.Flags().Parse(args); err != nil {
		t.Fatal(err)
	}

	// Assert
	if c.flagName != "nitro-dev" {
		t.Errorf("expected flag %q to be %q; got %q instead", "name", "nitro-dev", c.flagName)
	}
	if c.flagCpus != 2 {
		t.Errorf("expected flag %q to be %q; got %q instead", "cpus", 2, c.flagCpus)
	}
	if c.flagMemory != "2GB" {
		t.Errorf("expected flag %q to be %q; got %q instead", "memory", "2GB", c.flagMemory)
	}
	if c.flagDisk != "20GB" {
		t.Errorf("expected flag %q to be %q; got %q instead", "disk", "20GB", c.flagDisk)
	}
	if c.flagSkipInstall != true {
		t.Errorf("expected flag %q to be %v; got %v instead", "skip-install", true, c.flagSkipInstall)
	}
}

func TestInitCommand_Run(t *testing.T) {
	// create an empty config file
	configWithOutFile := viper.New()

	// create a config file with an example
	configWithFile := viper.New()
	configWithFile.SetConfigType("yaml")
	configWithFile.SetConfigFile("./testdata/nitro.yaml")
	_ = configWithFile.ReadInConfig()

	tests := []struct {
		name            string
		args            []string
		expected        []string
		chainedCommands []string
		want            int
		config          *viper.Viper
	}{
		{
			name:     "uses the flag arguments over config file or defaults",
			args:     []string{"-name", "example-test", "-cpu", "4"},
			expected: []string{"launch", "--name", "example-test", "--cpus", "4", "--mem", "2G", "--disk", "20G", "--cloud-init", "-"},
			want:     0,
			config:   configWithOutFile,
		},
		{
			name:     "uses the default if no flags are specified",
			args:     nil,
			expected: []string{"launch", "--name", "nitro-dev", "--cpus", "2", "--mem", "2G", "--disk", "20G", "--cloud-init", "-"},
			want:     0,
			config:   configWithOutFile,
		},
		//{
		//	name:            "will use the configuration file if provided",
		//	args:            nil,
		//	expected:        []string{"launch", "--name", "from-config-file", "--cpus", "4", "--mem", "4G", "--disk", "40G", "--cloud-init", "-"},
		//	want:            0,
		//	config:          configWithFile,
		//	chainedCommands: []string{"exec", "from-config-file", "--", "docker", "run", "-d", "--restart=always", "mysql:5.6", "-p", "3306:3306", "-e", "MYSQL_ROOT_PASSWORD=nitro", "-e", "MYSQL_DATABASE=nitro", "-e", "MYSQL_USER=nitro", "-e", "MYSQL_PASSWORD=nitro",},
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spyRunner := &SpyRunner{}

			_, c := testInitCommand(t, spyRunner, tt.config)

			if got := c.Run(tt.args); got != tt.want {
				t.Errorf("Run() = %v, want %v", got, tt.want)
			}

			if tt.expected != nil {
				if !reflect.DeepEqual(tt.expected, spyRunner.calls) {
					t.Errorf("wanted: \n%v \ngot: \n%v", tt.expected, spyRunner.calls)
				}
			}

			if tt.chainedCommands != nil {
				if !reflect.DeepEqual(tt.chainedCommands, spyRunner.chainedCalls) {
					t.Errorf("wanted: \n%v \ngot: \n%v", tt.chainedCommands, spyRunner.chainedCalls)
				}
			}
		})
	}
}
