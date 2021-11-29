package dockerbuild

import (
	"io"
	"os/exec"
)

var File = "docker"

func Build(r io.Reader, w io.Writer, path, image string) error {
	p, err := exec.LookPath(File)
	if err != nil {
		return err
	}

	c := exec.Command(p, "build", path, "--tag="+image)

	c.Stdin = r
	c.Stderr = w
	c.Stdout = w

	return c.Run()
}
