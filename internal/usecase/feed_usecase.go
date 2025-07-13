package usecase

import (
	"github.com/qrave1/gecko-eats/internal/infrastructure/postgres"
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
		&postgres.Feed{
			GeckoID:  geckoID,
			Date:     date,
			FoodType: foodType,
		},
	)
}

func (u *FeedUsecase) GetByGeckoID(geckoID string) ([]*postgres.Feed, error) {
	return u.repo.FeedsByGeckoID(geckoID, 100)
}

func (u *FeedUsecase) GetByDate(date string) ([]*postgres.Feed, error) {
	return u.repo.FeedsByDate(date)
}

func (u *FeedUsecase) DeleteAll(geckoID string) error {
	return u.repo.ClearFeed(geckoID)
}
