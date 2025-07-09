package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/qrave1/gecko-eats/cmd/config"
	"github.com/qrave1/gecko-eats/internal/infrastructure/telegram"
	"github.com/qrave1/gecko-eats/internal/repository"
	"github.com/qrave1/gecko-eats/internal/usecase"

	_ "github.com/lib/pq"
	tele "gopkg.in/telebot.v4"
)

// App represents the main application with all its dependencies
type App struct {
	Config        *config.Config
	Logger        *slog.Logger
	DB            *sqlx.DB
	Repository    repository.Repository
	Bot           *tele.Bot
	NotifyUsecase usecase.NotifyUsecaseInterface

	BotServer *telegram.BotServer
}

// NewApp creates a new application instance with all dependencies
func NewApp() (*App, error) {
	// Create config
	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	// Create logger
	logger := createLogger()

	// Create database connection
	db, err := createPostgresConnection(cfg)
	if err != nil {
		return nil, err
	}

	// Create repository
	repo := repository.NewSqlxRepository(db)

	geckoUsecase := usecase.NewGeckoUsecase(repo)

	feedUsecase := usecase.NewFeedUsecase(repo)

	// Create Telegram bot
	bot, err := createTelegramBot(cfg)
	if err != nil {
		return nil, err
	}

	// Create bot server
	botServer, err := telegram.NewBotServer(bot, geckoUsecase, feedUsecase, cfg.Bot.Whitelist)
	if err != nil {
		return nil, err
	}

	// Create notifyUsecase
	notifyUsecase := usecase.NewNotifyUsecase(bot, repo)

	return &App{
		Config:        cfg,
		Logger:        logger,
		DB:            db,
		Repository:    repo,
		Bot:           bot,
		BotServer:     botServer,
		NotifyUsecase: notifyUsecase,
	}, nil
}

// StartBot starts the Telegram bot server
func (a *App) StartBot(ctx context.Context) {
	a.BotServer.Start(ctx)
}

// RunNotifier runs the notification process once
func (a *App) RunNotifier(ctx context.Context) error {
	return a.NotifyUsecase.Notify(ctx, a.Config.Bot.NotifyUsers)
}

func createLogger() *slog.Logger {
	level := slog.LevelInfo

	replaceAttrFunc := func(groups []string, attributes slog.Attr) slog.Attr {
		if attributes.Key == "msg" {
			return slog.Attr{
				Key:   "message",
				Value: attributes.Value,
			}
		}

		return attributes
	}

	options := &slog.HandlerOptions{
		Level:       level,
		ReplaceAttr: replaceAttrFunc,
	}

	jsonHandler := slog.NewJSONHandler(os.Stdout, options)
	logger := slog.New(jsonHandler)
	slog.SetDefault(logger)

	return logger
}

func createPostgresConnection(cfg *config.Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSL,
	)

	conn, err := sqlx.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	// Check connection
	if err = conn.Ping(); err != nil {
		return nil, err
	}

	slog.Info(
		"open Postgres connection",
		slog.Any("DSN", dsn),
	)

	return conn, nil
}

func createTelegramBot(cfg *config.Config) (*tele.Bot, error) {
	return tele.NewBot(
		tele.Settings{
			Token:   cfg.Bot.Token,
			Poller:  &tele.LongPoller{Timeout: 5 * time.Second},
			Verbose: false,
		},
	)
}
