package app

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/qrave1/gecko-eats/internal/config"
	"github.com/qrave1/gecko-eats/internal/infrastructure/telegram"
	"github.com/qrave1/gecko-eats/internal/repository"
	"github.com/qrave1/gecko-eats/internal/usecase"

	_ "github.com/glebarez/sqlite"
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
func NewApp(configPath string) (*App, error) {
	// Create config
	cfg, err := config.New(configPath)
	if err != nil {
		return nil, err
	}

	// Create logger
	logger := createLogger()

	// Log startup
	logger.Info("starting application with config", slog.Any("config_path", configPath))

	// Create database connection
	db, err := createDBConnection(cfg)
	if err != nil {
		return nil, err
	}

	// Create repository
	repo := repository.NewSqlxRepository(db, cfg.Bot.Whitelist)

	// Create Telegram bot
	bot, err := createTelegramBot(cfg)
	if err != nil {
		return nil, err
	}

	// Create bot server
	botServer, err := telegram.NewBotServer(bot, repo, cfg.Bot.Whitelist)
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
	return a.NotifyUsecase.Notify(ctx, a.Config.Bot.NotifyUserIDs)
}

// Helper functions

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

func createDBConnection(cfg *config.Config) (*sqlx.DB, error) {
	conn, err := sqlx.Open("sqlite", cfg.Database.Path)
	if err != nil {
		return nil, err
	}

	// Create tables if they don't exist
	conn.Exec(`
	CREATE TABLE IF NOT EXISTS geckos
	(
		id         VARCHAR(36)  PRIMARY KEY,
		name       VARCHAR(255) NOT NULL UNIQUE,
		food_cycle VARCHAR(255)
	);
	
	CREATE TABLE IF NOT EXISTS feeds
	(
		date      VARCHAR(10)  NOT NULL,
		gecko_id    VARCHAR(36)  NOT NULL,
		food_type VARCHAR(255) NOT NULL,
		PRIMARY KEY (gecko_id, date)
	);
	`)

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
