package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/internal/client"
	"github.com/craftcms/nitro/internal/config"
	"github.com/craftcms/nitro/internal/nitro"
	"github.com/craftcms/nitro/internal/nitrod"
)

// TODO work this into a new commands pkg with a func that returns the command for testing.
// func NewInisetMaxExecutionCommand(addr IPer) &cobra.Command {
//	return &cobra.Command{
//		Use:   "max_execution_time",
//		Short: "Change max_execution_time",
//		Args:  cobra.MinimumNArgs(1),
//		RunE: func(cmd *cobra.Command, args []string) error {
// 			ip := addr.Get(flagMachineName)
// 			c := client.NewClient(ip, "50051")
//			php := config.GetString("php", flagPhpVersion)
//			resp, err := c.PhpIniSettings(cmd.Context(), &nitrod.ChangePhpIniSettingRequest{
//				Version: php,
//				Setting: nitrod.PhpIniSetting_MAX_EXECUTION_TIME,
//				Value:   args[0],
//			})
//			if err != nil {
//				return err
//			}
//
//			fmt.Println(resp.Message)
//
//			return nil
//		},
// }
var inisetMaxExecutionTimeCommand = &cobra.Command{
	Use:   "max_execution_time",
	Short: "Change max_execution_time",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName
		runner := nitro.NewMultipassRunner("multipass")
		ip := nitro.IP(machine, runner)
		c, err := client.NewClient(ip, "50051")
		if err != nil {
			return err
		}
		php := config.GetString("php", flagPhpVersion)

		resp, err := c.PhpIniSettings(cmd.Context(), &nitrod.ChangePhpIniSettingRequest{
			Version: php,
			Setting: nitrod.PhpIniSetting_MAX_EXECUTION_TIME,
			Value:   args[0],
		})
		if err != nil {
			return err
		}

		if !flagSilent {
			fmt.Println(resp.Message)
		}

		return nil
	},
}

func init() {
	inisetMaxExecutionTimeCommand.Flags().BoolVar(&flagSilent, "silent", false, "Run command with no output")
}
