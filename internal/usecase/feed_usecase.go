package usecase

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/qrave1/gecko-eats/internal/domain"
	"github.com/qrave1/gecko-eats/internal/domain/input"
	"github.com/qrave1/gecko-eats/internal/repository"
)

type FeedUsecase struct {
	repo repository.Repository
}

func NewFeedUsecase(repo repository.Repository) *FeedUsecase {
	return &FeedUsecase{
		repo: repo,
	}
}

func (u *FeedUsecase) Create(geckoID, date, foodType string) error {
	return u.repo.AddFeed(
		&domain.Feed{
			GeckoID:  geckoID,
			Date:     date,
			FoodType: foodType,
		},
	)
}

func (u *FeedUsecase) CreateSchedule(
	gecko *domain.Gecko,
	createFeedingsInput *input.CreateFeedingsInput,
) (int, error) {
	var cycle []string

	err := json.Unmarshal([]byte(gecko.FoodCycle), &cycle)
	if err != nil {
		return 0, err
	}

	createdFeedsCount := 0

	endDate := createFeedingsInput.StartDate.AddDate(0, 1, 0) // +1 месяц

	for d := createFeedingsInput.StartDate; d.Before(endDate); d = d.AddDate(0, 0, createFeedingsInput.Interval) {
		event := time.Date(d.Year(), d.Month(), d.Day(), 18, 0, 0, 0, d.Location())

		err := u.Create(gecko.ID, event.Format("2006-01-02"), cycle[createFeedingsInput.StartIx])

		if err != nil {
			slog.Error(
				"add feed",
				slog.Any("error", err),
				slog.String("gecko_id", gecko.ID),
				slog.String("date", event.Format("2006-01-02")),
				slog.String("food_type", cycle[createFeedingsInput.StartIx]),
			)
		}

		createdFeedsCount++

		// счётчик для следующего типа еды
		if createFeedingsInput.StartIx < len(cycle)-1 {
			createFeedingsInput.StartIx++
		} else {
			createFeedingsInput.StartIx = 0
		}
	}

	return createdFeedsCount, nil
}

func (u *FeedUsecase) GetByGeckoID(geckoID string) ([]*domain.Feed, error) {
	return u.repo.FeedsByGeckoID(geckoID, 100)
}

func (u *FeedUsecase) GetByDate(date string) ([]*domain.Feed, error) {
	return u.repo.FeedsByDate(date)
}

func (u *FeedUsecase) DeleteAll(geckoID string) error {
	return u.repo.ClearFeed(geckoID)
}
