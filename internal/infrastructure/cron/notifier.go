package cron

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

type Notifier struct {
	bot     *tele.Bot
	repo    repository.Repository
	userIDs []int64
}

func NewNotifier(bot *tele.Bot, repo repository.Repository, userIDs []int64) *Notifier {
	return &Notifier{userIDs: userIDs, repo: repo, bot: bot}
}

func (c *Notifier) Notify(_ context.Context) error {
	date := time.Now().UTC().Format("2006-01-02")

	var messageBuilder strings.Builder

	messageBuilder.WriteString(fmt.Sprintf("📅 Кормление на %s:\n", date))

	feedings, err := c.repo.FeedingsByDate(date)

	if err != nil {
		return fmt.Errorf("get feedings: %w", err)
	}

	pets, err := c.repo.Pets()

	if err != nil {
		return fmt.Errorf("get pets: %w", err)
	}

	for _, feeding := range feedings {
		var pet *sql.Pet

		for _, p := range pets {
			if p.ID == feeding.PetID {
				pet = &p

				break
			}
		}

		if pet == nil {
			slog.Error(
				"pet not found for feeding ID",
				slog.String("feeding_id", feeding.PetID),
			)

			continue
		}

		messageBuilder.WriteString(fmt.Sprintf("🦎 %s — корм: %s\n", pet.Name, feeding.FoodType))
	}

	for _, uid := range c.userIDs {
		_, err = c.bot.Send(
			&tele.User{ID: uid},
			messageBuilder.String(),
		)
	}

	return nil
}
