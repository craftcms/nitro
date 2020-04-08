package cmd

import (
	"github.com/craftcms/nitro/internal/action"
)

func Run(r ShellRunner, actions []action.Action) error {
	for _, a := range actions {
		// if this is the launch action, check for input
		if a.Type == "launch" && a.Input != "" {
			if err := r.SetInput(a.Input); err != nil {
				return err
			}
		}

		r.UseSyscall(a.UseSyscall)

		// only return an error if its not nil
		if err := r.Run(a.Args); err != nil {
			return err
		}
	}

	return nil
}
