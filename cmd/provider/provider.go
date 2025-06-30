package provider

import (
	"log/slog"
	"os"
	"time"

	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
	"github.com/qrave1/gecko-eats/internal/config"
	"github.com/qrave1/gecko-eats/internal/infrastructure/cron"
	"github.com/qrave1/gecko-eats/internal/infrastructure/telegram"
	"github.com/qrave1/gecko-eats/internal/repository"

	_ "github.com/glebarez/sqlite"

	tele "gopkg.in/telebot.v4"
)

var BotServiceSet = wire.NewSet(
	ProvideConfig,
	ProvideLogger,
	ProvideSQLXConnection,
	ProvideRepository,
	ProvideTeleBot,
	ProvideBotServer,
)

var NotifyServiceSet = wire.NewSet(
	ProvideConfig,
	ProvideLogger,
	ProvideSQLXConnection,
	ProvideRepository,
	ProvideTeleBot,
	ProvideNotifier,
)

func ProvideConfig() (*config.Config, error) {
	return config.New(CONFIG_PATH)
}

func ProvideLogger(_ *config.Config) (*slog.Logger, error) {
	level := slog.LevelInfo

	// TODO: log level from config

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

	//handlers := make([]slog.Handler, 0, 2)

	jsonHandler := slog.NewJSONHandler(os.Stdout, options)

	//handlers = append(handlers, jsonHandler)

	logger := slog.New(jsonHandler)

	slog.SetDefault(logger)

	logger.Info("start with config", slog.Any("config_path", CONFIG_PATH))

	return logger, nil
}

func ProvideSQLXConnection(cfg *config.Config) (*sqlx.DB, error) {
	conn, err := sqlx.Open("sqlite", cfg.Database.Path)

	if err != nil {
		return nil, err
	}

	conn.Exec(
		`
	CREATE TABLE IF NOT EXISTS pets
	(
		id         VARCHAR(36)  PRIMARY KEY,
		name       VARCHAR(255) NOT NULL UNIQUE,
		food_cycle VARCHAR(255)
	);
	
	CREATE TABLE IF NOT EXISTS feedings
	(
		id        INTEGER PRIMARY KEY AUTOINCREMENT,
		date      VARCHAR(10)  NOT NULL UNIQUE,
		pet_id    VARCHAR(36)  NOT NULL,
		food_type VARCHAR(255) NOT NULL
	);
`,
	)

	return conn, nil
}

func ProvideRepository(cfg *config.Config, db *sqlx.DB) repository.Repository {
	return repository.NewSqlxRepository(db, cfg.Bot.Whitelist)
}

func ProvideTeleBot(cfg *config.Config) (*tele.Bot, error) {
	bot, err := tele.NewBot(
		tele.Settings{
			Token:  cfg.Bot.Token,
			Poller: &tele.LongPoller{Timeout: 5 * time.Second},
		},
	)

	if err != nil {
		return nil, err
	}

	return bot, nil
}

func ProvideBotServer(
	_ *slog.Logger,
	bot *tele.Bot,
	repo repository.Repository,
	cfg *config.Config,
) (*telegram.BotServer, error) {
	return telegram.NewBotServer(bot, repo, cfg.Bot.Whitelist)
}

func ProvideNotifier(
	_ *slog.Logger,
	bot *tele.Bot,
	repo repository.Repository,
	cfg *config.Config,
) *cron.Notifier {
	return cron.NewNotifier(bot, repo, cfg.Bot.NotifyUserIDs)
}
