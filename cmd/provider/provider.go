package provider

import (
	"time"

	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
	"github.com/qrave1/gecko-eats/internal/config"
	"github.com/qrave1/gecko-eats/internal/infrastructure/cron"
	"github.com/qrave1/gecko-eats/internal/infrastructure/telegram"
	"github.com/qrave1/gecko-eats/internal/repository"
	tele "gopkg.in/telebot.v4"
)

var BotServiceSet = wire.NewSet(
	ProvideConfig,
	ProvideSQLXConnection,
	ProvideRepository,
	ProvideTeleBot,
	ProvideBotServer,
)

var NotifyServiceSet = wire.NewSet(
	ProvideConfig,
	ProvideSQLXConnection,
	ProvideRepository,
	ProvideTeleBot,
	ProvideNotifier,
)

func ProvideConfig() (*config.Config, error) {
	// TODO: delete hardcode
	path := "./config.yaml"

	return config.New(path)
}

func ProvideSQLXConnection(cfg *config.Config) (*sqlx.DB, error) {
	conn, err := sqlx.Open("sqlite", cfg.Database.Path)

	if err != nil {
		return nil, err
	}

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
	bot *tele.Bot,
	repo repository.Repository,
	cfg *config.Config,
) (*telegram.BotServer, error) {
	return telegram.NewBotServer(bot, repo, cfg.Bot.Whitelist)
}

func ProvideNotifier(
	bot *tele.Bot,
	repo repository.Repository,
	cfg *config.Config,
) *cron.Notifier {
	return cron.NewNotifier(bot, repo, cfg.Bot.NotifyUserIDs)
}
