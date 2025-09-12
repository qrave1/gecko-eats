package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/qrave1/gecko-eats/internal/domain"
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

	feeds, err := n.repo.FeedsByDate(date)
	if err != nil {
		return fmt.Errorf("get feeds: %w", err)
	}

	if len(feeds) == 0 {
		slog.Info("no feed found")

		return nil
	}

	geckos, err := n.repo.Geckos()
	if err != nil {
		return fmt.Errorf("get geckos: %w", err)
	}

	var messageBuilder strings.Builder
	messageBuilder.WriteString(fmt.Sprintf("📅 Кормление на %s:\n", date))

	for _, feed := range feeds {
		var gecko *domain.Gecko

		for _, p := range geckos {
			if p.ID == feed.GeckoID {
				gecko = p

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

	notifyMessage := messageBuilder.String()

	for _, uid := range userIDs {
		_, err = n.bot.Send(
			&tele.User{ID: uid},
			notifyMessage,
		)
	}

	return nil
}
