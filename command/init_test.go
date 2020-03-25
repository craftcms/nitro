package command

import (
	"reflect"
	"strings"
	"testing"
)

func TestInitCommand_Synopsis(t *testing.T) {
	// Arrange
	c := InitCommand{}
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
	c := InitCommand{}
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
	c := InitCommand{}
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

type RunnerSpy struct {
	Calls []string
}

func (r *RunnerSpy) Run(args []string) error {
	r.Calls = args
	return nil
}

func TestInitCommand_Run(t *testing.T) {
	spyRunner := &RunnerSpy{}
	type fields struct {
		runner ShellRunner
	}
	tests := []struct {
		name     string
		fields   fields
		args     []string
		expected []string
		want     int
	}{
		{
			name:     "Run gets the right arguments passed to multipass",
			fields:   fields{runner: spyRunner},
			args:     []string{"-name=example-test", "-cpu=2"},
			expected: []string{"multipass", "launch", "--name", "example-test", "--cpus", "2"},
			want:     0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &InitCommand{
				runner: tt.fields.runner,
			}

			if got := c.Run(tt.args); got != tt.want {
				t.Errorf("Run() = %v, want %v", got, tt.want)
			}

			if !reflect.DeepEqual(tt.expected, spyRunner.Calls) {
				t.Errorf("wanted calls %v got %v", tt.expected, spyRunner.Calls)
			}
		})
	}
}
