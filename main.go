package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/qrave1/gecko-eats/cmd/provider"
	"github.com/qrave1/gecko-eats/cmd/wire"
	"github.com/qrave1/gecko-eats/internal/infrastructure/sql"
	"github.com/urfave/cli/v3"
)

func main() {
	ctx := context.Background()

	app := &cli.Command{
		Name:  "gecko-feeder",
		Usage: "bot for feeding geckos",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "config.yaml",
				Usage:   "path to YAML config file",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "notify",
				Usage: "cron job to notify about today's feedings",
				Action: func(ctx context.Context, c *cli.Command) error {
					provider.CONFIG_PATH = c.String("config")

					_, err := wire.InitializeNotifyService(ctx)

					if err != nil {
						panic(err)
					}

					return nil
				},
			},
			{
				Name:  "migrate",
				Usage: "run SQL database migrations",
				Action: func(ctx context.Context, c *cli.Command) error {
					provider.CONFIG_PATH = c.String("config")

					return sql.RunMigrations(provider.CONFIG_PATH)
				},
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			provider.CONFIG_PATH = c.String("config")

			_, err := wire.InitializeBotService(ctx)

			if err != nil {
				panic(err)
			}

			exit := make(chan os.Signal, 1)

			signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-exit:
			}

			return nil
		},
	}

	if err := app.Run(ctx, os.Args); err != nil {
		panic(err)
	}
}
