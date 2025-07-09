package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/qrave1/gecko-eats/cmd/app"
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
					application, err := app.NewApp()
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
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			// Initialize the application
			application, err := app.NewApp()
			if err != nil {
				panic(err)
			}

			// Start the bot
			application.StartBot(ctx)

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

	if err := cliApp.Run(ctx, os.Args); err != nil {
		panic(err)
	}
}
