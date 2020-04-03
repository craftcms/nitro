package nitro

import (
	"log"

	"github.com/craftcms/nitro/internal/command"
)

func Run(runner command.Runner, commands []Command) error {
	for _, c := range commands {
		if c.Input != "" {
			if err := runner.SetInput(c.Input); err != nil {
				return err
			}
		}

		if c.Chainable == false {
			runner.UseSyscall(true)
		}

		var preArgs []string
		switch c.Type {
		case "launch":
			preArgs = append(preArgs, "launch")
			preArgs = append(preArgs, c.Args...)
		case "shell":
			runner.UseSyscall(true)
			preArgs = append(preArgs, "shell", c.Machine)
		default:
			preArgs = append(preArgs, "exec")
			preArgs = append(preArgs, c.Args...)
		}

		if err := runner.Run(preArgs); err != nil {
			log.Println("error in runner.Run:", err.Error())
			log.Println(preArgs)
			return err
		}
	}

	return nil
}
