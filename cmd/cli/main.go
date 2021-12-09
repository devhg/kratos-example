package main

import (
	"log"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	"github.com/urfave/cli/v2"

	"github.com/devhg/kratos-example/pkg/migrate"
)

func main() {
	app := cli.NewApp()
	app.Usage = "AppCli工具"

	app.Commands = []*cli.Command{
		{
			Name:        "migrate",
			Usage:       "同步数据库",
			Description: "go run xxx/main.go migrate",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "config",
					Usage:       "config path, eg: -conf config.yaml",
					DefaultText: "./configs",
					Value:       "./configs",
				}},
			Action: func(ctx *cli.Context) error {
				conf, err := filepath.Abs(ctx.String("config"))
				if err != nil {
					return err
				}
				return migrate.Migrate(conf, nil)
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
