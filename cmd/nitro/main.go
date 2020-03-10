package main

import (
	"context"
	"log"
	"os"
	"os/exec"

	"github.com/craftcms/nitro/internal/app"
	"github.com/craftcms/nitro/internal/executor"
	"github.com/craftcms/nitro/internal/runner"
)

func run(args []string) {
	// find the path to multipass and set value in context
	multipass, err := exec.LookPath("multipass")
	if err != nil {
		log.Fatal(err)
	}

	run := runner.NewRunner("multipass")

	ctx := context.WithValue(context.Background(), "multipass", multipass)

	if err := app.NewApp(executor.NewSyscallExecutor("multipass"), run).RunContext(ctx, args); err != nil {
		log.Fatal(err)
	}
}

func main() {
	run(os.Args)
}
