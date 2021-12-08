package main

import (
	"context"
	"entgo.io/ent/entc/integration/ent/migrate"
	"github.com/devhg/kratos-example/internal/data/ent"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	// "github.com/devhg/kratos-example/app/migrate"

	_ "github.com/go-sql-driver/mysql"
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
				// conf, err := filepath.Abs(ctx.String("config"))
				// if err != nil {
				// 	return err
				// }
				// return migrate.Migrate(conf, nil)
				dataSource := "root:root@tcp(127.0.0.1:3306)/testdb?parseTime=True"
				client, err := ent.Open("mysql", dataSource)
				if err != nil {
					log.Fatalf("failed creating schema resources: %v", err)
				}
				defer client.Close()

				if err := client.Debug().Schema.Create(
					context.Background(),
					migrate.WithDropIndex(true),
					migrate.WithDropColumn(true)); err != nil {
					log.Fatalf("failed creating schema resources: %v", err)
				}
				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
