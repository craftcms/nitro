package command

import (
	"github.com/urfave/cli/v2"
)

var bootstrapFlag = &cli.BoolFlag{
	Name:        "bootstrap",
	Usage:       "Bootstrap the machine with defaults",
	Value:       true,
	DefaultText: "true",
}

var cpusFlag = &cli.Int64Flag{
	Name:        "cpus",
	Usage:       "The number of CPUs to assign the machine",
	Required:    false,
	Value:       2,
	DefaultText: "2",
}

var databaseFlag = &cli.StringFlag{
	Name:        "database",
	Usage:       "Provide version of PHP",
	Value:       "mariadb",
	DefaultText: "mariadb",
}

var diskFlag = &cli.StringFlag{
	Name:        "disk",
	Usage:       "The amount of disk to assign the machine",
	Required:    false,
	Value:       "5G",
	DefaultText: "5G",
}

var memoryFlag = &cli.StringFlag{
	Name:        "memory",
	Usage:       "The amount of memory to assign the machine",
	Required:    false,
	Value:       "2G",
	DefaultText: "2G",
}

var permanentDeleteFlag = &cli.BoolFlag{
	Name:        "permanent",
	Usage:       "Permanently delete a machine",
	Value:       false,
	DefaultText: "false",
}

var phpVersionFlag = &cli.StringFlag{
	Name:        "php-version",
	Usage:       "Provide version of PHP",
	Value:       "7.4",
	DefaultText: "7.4",
}

var postgresFlag = &cli.BoolFlag{
	Name:        "postgres",
	Usage:       "Enter a postgres shell",
	Value:       false,
	DefaultText: "false",
}

var pathFlag = &cli.StringFlag{
	Name:        "path",
	Usage:       "The path to the site",
	Value:       "",
}

var publicDirFlag = &cli.StringFlag{
	Name:        "public-dir",
	Usage:       "The public directory for the server",
	Value:       "web",
	DefaultText: "web",
}

var serviceMySqlFlag = &cli.BoolFlag{
	Name:        "mysql",
	Usage:       "Affect MySQL service",
	Value:       false,
	DefaultText: "false",
}

var servicePostgresFlag = &cli.BoolFlag{
	Name:        "postgres",
	Usage:       "Affect PostgreSQL service",
	Value:       false,
	DefaultText: "false",
}

var serviceNginxFlag = &cli.BoolFlag{
	Name:        "nginx",
	Usage:       "Affect Nginx service",
	Value:       false,
	DefaultText: "false",
}

var serviceRedisFlag = &cli.BoolFlag{
	Name:        "redis",
	Usage:       "Affect redis service",
	Value:       false,
	DefaultText: "false",
}
