//go:build wireinject
// +build wireinject

package wire

import (
	"context"

	"github.com/google/wire"
	"github.com/qrave1/gecko-eats/cmd/provider"
)

func InitializeBotService(_ context.Context) (*provider.BotService, error) {
	wire.Build(
		provider.BotServiceSet,
		provider.ProvideBotService,
	)

	return &provider.BotService{}, nil
}

func InitializeNotifyService(_ context.Context) (*provider.NotifyService, error) {
	wire.Build(
		provider.NotifyServiceSet,
		provider.ProvideNotifyService,
	)

	return &provider.NotifyService{}, nil
}
