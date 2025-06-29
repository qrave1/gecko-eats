package provider

import (
	"context"

	"github.com/qrave1/gecko-eats/internal/infrastructure/cron"
	"github.com/qrave1/gecko-eats/internal/infrastructure/telegram"
)

type BotService struct {
	BotServer *telegram.BotServer
}

func ProvideBotService(ctx context.Context, botServer *telegram.BotServer) *BotService {
	bs := &BotService{
		BotServer: botServer,
	}

	bs.BotServer.Start(ctx)

	return bs
}

type NotifyService struct {
	Notifier *cron.Notifier
}

func ProvideNotifyService(notifier *cron.Notifier) *NotifyService {
	return &NotifyService{
		Notifier: notifier,
	}
}
