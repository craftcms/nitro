package nitro

import (
	"log"
)

func Run(runner ShellRunner, commands []Command) error {
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
		case "mount":
			preArgs = append(preArgs, "mount")
			preArgs = append(preArgs, c.Args...)
		case "start":
			preArgs = append(preArgs, "start")
			preArgs = append(preArgs, c.Args...)
		case "stop":
			preArgs = append(preArgs, "stop")
			preArgs = append(preArgs, c.Args...)
		case "info":
			preArgs = append(preArgs, "info")
			preArgs = append(preArgs, c.Args...)
		case "delete":
			preArgs = append(preArgs, "delete")
			preArgs = append(preArgs, c.Args...)
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
			log.Println("error in nitro runner:", err.Error())
			return err
		}
	}

	return nil
}
