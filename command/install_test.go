package command

import (
	"reflect"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/spf13/viper"
)

func testInstallCommand(t testing.TB, runner *SpyRunner, v *viper.Viper) (*cli.MockUi, *InstallCommand) {
	t.Helper()

	ui := cli.NewMockUi()
	coreCmd := &CoreCommand{
		UI:     ui,
		Runner: runner,
		Config: v,
	}

	return ui, &InstallCommand{CoreCommand: coreCmd}
}

func TestInstallCommand_Synopsis(t *testing.T) {
	// Arrange
	v := viper.New()
	_, c := testInstallCommand(t, nil, v)
	expected := "install software on machine"

	// Act
	actual := c.Synopsis()

	// Assert
	if actual != expected {
		t.Errorf("expected %q; got %q instead", expected, actual)
	}
}

func TestInstallCommand_Help(t *testing.T) {
	// Arrange
	v := viper.New()
	_, c := testInstallCommand(t, nil, v)
	expected := strings.TrimSpace(`
Usage: nitro install [options]
  This command install software on a virtual machine.
  
  Install software on a virtual machine and override the software version:
      $ nitro install -php-version=7.4 -database-engine=mysql -database-version=5.7
  
  Install software on a virtual machine with the default options:
      $ nitro install

  The default options will be the latest version of PHP (version 7.4) and MySQL 5.7.
`)

	// Act
	actual := c.Help()

	// Assert
	if actual != expected {
		t.Errorf("expected %q; got %q instead", expected, actual)
	}
}

func TestInstallCommand_Flags(t *testing.T) {
	// Arrange
	v := viper.New()
	v.SetConfigFile("./testdata/nitro.yaml")
	v.SetConfigType("yaml")
	_ = v.ReadInConfig()
	_, c := testInstallCommand(t, nil, v)
	args := []string{"-php-version=7.3", "-database-engine=mysql", "-database-version=5.7"}

	// Act
	if err := c.Flags().Parse(args); err != nil {
		t.Fatal(err)
	}

	// Assert
	if c.flagName != "from-config-file" {
		t.Errorf("expected flag %q to be %q; got %q instead", "name", "from-config-file", c.flagName)
	}
	if c.flagPhpVersion != "7.3" {
		t.Errorf("expected flag %q to be %q; got %q instead", "php-version", "7.4", c.flagPhpVersion)
	}
	if c.flagDatabaseEngine != "mysql" {
		t.Errorf("expected flag %q to be %q; got %q instead", "database-engine", "mysql", c.flagDatabaseEngine)
	}
}

func TestInstallCommand_FlagsUsesConfigFile(t *testing.T) {
	// Arrange
	v := viper.New()
	v.SetConfigFile("./testdata/nitro.yaml")
	v.SetConfigType("yaml")
	_ = v.ReadInConfig()
	_, c := testInstallCommand(t, nil, v)

	// Act
	if err := c.Flags().Parse(nil); err != nil {
		t.Fatal(err)
	}

	// Assert
	if c.flagName != "from-config-file" {
		t.Errorf("expected flag %q to be %q; got %q instead", "name", "from-config-file", c.flagName)
	}
	if c.flagPhpVersion != "7.3" {
		t.Errorf("expected flag %q to be %q; got %q instead", "php-version", "7.3", c.flagPhpVersion)
	}
	if c.flagDatabaseEngine != "mysql" {
		t.Errorf("expected flag %q to be %q; got %q instead", "database-engine", "mysql", c.flagDatabaseEngine)
	}
	if c.flagDatabaseVersion != "5.6" {
		t.Errorf("expected flag %q to be %q; got %q instead", "database-version", "5.6", c.flagDatabaseVersion)
	}
}

func TestInstallCommand_Run(t *testing.T) {
	configWithFile := viper.New()
	configWithFile.SetConfigFile("./testdata/nitro.yaml")
	configWithFile.SetConfigType("yaml")
	_ = configWithFile.ReadInConfig()
	spyRunner := &SpyRunner{}

	tests := []struct {
		name            string
		args            []string
		expected        []string
		chainedCommands []string
		config          *viper.Viper
		want            int
	}{
		{
			name:            "uses the flag arguments over configWithFile file or defaults",
			args:            nil,
			expected:        []string{"exec", "from-config-file", "--", "sudo", "apt", "install", "-y", "php7.3", "php7.3-mbstring", "php7.3-cli", "php7.3-curl", "php7.3-fpm", "php7.3-gd", "php7.3-intl", "php7.3-json", "php7.3-mysql", "php7.3-opcache", "php7.3-pgsql", "php7.3-zip", "php7.3-xml", "php-xdebug", "php-imagick"},
			chainedCommands: []string{"exec", "from-config-file", "--", "docker", "run", "-d", "--restart=always", "mysql:5.6", "-p", "3306:3306", "-e", "MYSQL_ROOT_PASSWORD=nitro", "-e", "MYSQL_DATABASE=nitro", "-e", "MYSQL_USER=nitro", "-e", "MYSQL_PASSWORD=nitro",},
			want:            0,
			config:          configWithFile,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, c := testInstallCommand(t, spyRunner, tt.config)
			if got := c.Run(tt.args); got != tt.want {
				t.Errorf("Run() = %configWithFile, want %configWithFile", got, tt.want)
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
