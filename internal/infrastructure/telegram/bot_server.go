package telegram

import (
	"context"
	"log/slog"

	"github.com/qrave1/gecko-eats/internal/infrastructure/telegram/middleware"
	"github.com/qrave1/gecko-eats/internal/repository"
	tele "gopkg.in/telebot.v4"
)

type BotServer struct {
	bot  *tele.Bot
	repo repository.Repository
}

func NewBotServer(bot *tele.Bot, repo repository.Repository, whitelist []int64) (*BotServer, error) {
	botServer := &BotServer{
		bot:  bot,
		repo: repo,
	}

	botServer.bot.Use(middleware.WhitelistOnly(whitelist))

	return botServer, nil
}

func (b *BotServer) Start(ctx context.Context) {
	go b.bot.Start()

	select {
	case <-ctx.Done():
		b.bot.Stop()

		slog.Info("bot stopped")
	}
}
