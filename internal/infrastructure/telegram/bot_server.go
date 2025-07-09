package telegram

import (
	"context"
	"log/slog"

	"github.com/qrave1/gecko-eats/internal/infrastructure/telegram/middleware"
	"github.com/qrave1/gecko-eats/internal/usecase"
	tele "gopkg.in/telebot.v4"
	teleMid "gopkg.in/telebot.v4/middleware"
)

type BotServer struct {
	bot *tele.Bot

	geckoUsecase *usecase.GeckoUsecase
	feedUsecase  *usecase.FeedUsecase
}

func NewBotServer(
	bot *tele.Bot,
	geckoUsecase *usecase.GeckoUsecase,
	feedUsecase *usecase.FeedUsecase,
	whitelist []int64,
) (*BotServer, error) {
	botServer := &BotServer{
		bot:          bot,
		geckoUsecase: geckoUsecase,
		feedUsecase:  feedUsecase,
	}

	botServer.registerHandlers()

	botServer.bot.Use(middleware.WhitelistOnly(whitelist))
	botServer.bot.Use(teleMid.AutoRespond())

	return botServer, nil
}

func (b *BotServer) Start(ctx context.Context) {
	go b.bot.Start()

	slog.Info(
		"bot started",
		slog.Any("username", b.bot.Me.Username),
	)

	select {
	case <-ctx.Done():
		b.bot.Stop()

		slog.Info(
			"bot stopped",
			slog.Any("username", b.bot.Me.Username),
		)
	}
}
