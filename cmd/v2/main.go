package main

import (
	"os"

	"github.com/craftcms/nitro/pkg/cmd/destroy"
	"github.com/craftcms/nitro/pkg/cmd/initcmd"
	"github.com/craftcms/nitro/pkg/cmd/restart"
	"github.com/craftcms/nitro/pkg/cmd/start"
	"github.com/craftcms/nitro/pkg/cmd/stop"
	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use:          "nitro",
	Short:        "Local Craft CMS dev made easy",
	Long:         `Nitro is a command-line tool focused on making local Craft CMS development quick and easy.`,
	RunE:         rootMain,
	SilenceUsage: true,
}

func rootMain(command *cobra.Command, _ []string) error {
	return command.Help()
}

func init() {
	flags := rootCommand.PersistentFlags()

	// set global flags
	flags.String("environment", "nitro-dev", "The name of the environment")

	// register all of the commands
	commands := []*cobra.Command{
		initcmd.InitCommand,
		stop.StopCommand,
		start.StartCommand,
		destroy.DestroyCommand,
		restart.RestartCommand,
	}

	rootCommand.AddCommand(commands...)
}

func main() {
	// Execute the root command.
	if err := rootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}

// package main

// import (
// 	"context"
// 	"flag"
// 	"fmt"
// 	"log"
// 	"os"
// 	"path/filepath"

// 	"github.com/craftcms/nitro/pkg/client"
// )

// func main() {
// 	name := flag.String("machine", "nitro-dev", "the name of the machine")
// 	stop := flag.Bool("stop", false, "stop the containers")
// 	install := flag.String("install", "", "the path to the craft install")
// 	flag.Parse()

// 	ctx := context.Background()

// 	cli, err := client.NewClient()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	if *install != "" {
// 		abs, err := filepath.Abs(*install)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		if err := cli.Composer(ctx, abs, "1", "install"); err != nil {
// 			log.Fatal(err)
// 		}

// 		os.Exit(0)
// 	}

// 	if *stop {
// 		if err := cli.Stop(ctx, *name, os.Args); err != nil {
// 			fmt.Println("Error:", err)
// 			os.Exit(1)
// 		}
// 	} else {
// 		if err := cli.Init(ctx, *name, os.Args); err != nil {
// 			fmt.Println("Error:", err)
// 			os.Exit(1)
// 		}
// 	}
// }
