package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pressly/goose/v3"
	"github.com/qrave1/gecko-eats/cmd/application"
	"github.com/qrave1/gecko-eats/cmd/config"
	"github.com/qrave1/gecko-eats/internal/infrastructure/postgres"
	"github.com/urfave/cli/v3"
)

func main() {
	ctx := context.Background()

	cliApp := &cli.Command{
		Name:  "gecko-feeder",
		Usage: "bot for feed geckos",
		Commands: []*cli.Command{
			{
				Name:  "notify",
				Usage: "cron job to notify about today's feeds",
				Action: func(ctx context.Context, c *cli.Command) error {
					// Initialize the application
					application, err := application.NewApp(ctx)
					if err != nil {
						panic(err)
					}

					// Run the notifier
					if err = application.RunNotifier(ctx); err != nil {
						panic(err)
					}

					return nil
				},
			},
			{
				Name:  "migrate",
				Usage: "migrate database",
				Action: func(ctx context.Context, c *cli.Command) error {
					cfg, err := config.New()
					if err != nil {
						panic(err)
					}

					conn, err := application.NewPostgresConnection(ctx, cfg)
					if err != nil {
						panic(err)
					}

					goose.SetBaseFS(postgres.MigrationsEmbed)

					err = goose.SetDialect(string(goose.DialectPostgres))
					if err != nil {
						panic(err)
					}

					err = goose.RunContext(
						ctx,
						c.Args().First(),
						conn.DB,
						"migrations",
						c.Args().Slice()...,
					)
					if err != nil {
						panic(err)
					}

					return nil
				},
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			// Initialize the application
			app, err := application.NewApp(ctx)
			if err != nil {
				panic(err)
			}

			// Start the bot
			err = app.Start(ctx)

			if err != nil {
				panic(err)
			}

			exit := make(chan os.Signal, 1)
			signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
			<-exit

			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			if err = app.Shutdown(ctx); err != nil {
				panic(err)
			}

			return nil
		},
	}

	if err := cliApp.Run(ctx, os.Args); err != nil {
		panic(err)
	}
}
