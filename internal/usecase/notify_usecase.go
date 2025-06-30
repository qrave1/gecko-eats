package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/qrave1/gecko-eats/internal/infrastructure/sql"
	"github.com/qrave1/gecko-eats/internal/repository"
	tele "gopkg.in/telebot.v4"
)

type NotifyUsecaseInterface interface {
	Notify(ctx context.Context, userIDs []int64) error
}

type NotifyUsecase struct {
	bot  *tele.Bot
	repo repository.Repository
}

func NewNotifyUsecase(bot *tele.Bot, repo repository.Repository) *NotifyUsecase {
	return &NotifyUsecase{
		bot:  bot,
		repo: repo,
	}
}

func (n *NotifyUsecase) Notify(_ context.Context, userIDs []int64) error {
	date := time.Now().UTC().Format("2006-01-02")

	var messageBuilder strings.Builder

	messageBuilder.WriteString(fmt.Sprintf("📅 Кормление на %s:\n", date))

	feeds, err := n.repo.feedsByDate(date)

	if err != nil {
		return fmt.Errorf("get feeds: %w", err)
	}

	geckos, err := n.repo.Geckos()

	if err != nil {
		return fmt.Errorf("get geckos: %w", err)
	}

	for _, feed := range feeds {
		var gecko *sql.Gecko

		for _, p := range geckos {
			if p.ID == feed.GeckoID {
				gecko = &p

				break
			}
		}

		if gecko == nil {
			slog.Error(
				"gecko not found for feed ID",
				slog.String("feed_id", feed.GeckoID),
			)

			continue
		}

		messageBuilder.WriteString(fmt.Sprintf("🦎 %s — корм: %s\n", gecko.Name, feed.FoodType))
	}

	for _, uid := range userIDs {
		_, err = n.bot.Send(
			&tele.User{ID: uid},
			messageBuilder.String(),
		)
	}

	return nil
}
