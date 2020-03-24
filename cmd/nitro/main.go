package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"

	"github.com/craftcms/nitro/internal/app"
	"github.com/craftcms/nitro/internal/command"
)

func run(args []string) {
	if err := app.NewApp(command.NewRunner("multipass")).Run(args); err != nil {
		log.Fatal(err)
	}
}

func init() {
	v := viper.New()
	v.SetConfigName("nitro")
	v.SetConfigType("yml")
	v.AddConfigPath(".")

	// set defaults
	v.SetDefault("name", "nitro-dev")
	v.SetDefault("system.php", "7.4")
	v.SetDefault("system.cpus", "4")
	v.SetDefault("system.memory", "2G")
	v.SetDefault("system.disk", "20G")

	err := v.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func main() {

	args := os.Args[:1]
	log.Println("args:", args)
	args = append(args, "--machine="+viper.GetString("name"))
	log.Println("appended:", args)
	log.Println("after:", os.Args[len(os.Args)-1:len(os.Args)])
	//run(args)
}
