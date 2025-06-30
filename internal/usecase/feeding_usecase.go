package usecase

import (
	"github.com/qrave1/gecko-eats/internal/infrastructure/sql"
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
		&sql.Feed{
			GeckoID:  geckoID,
			Date:     date,
			FoodType: foodType,
		},
	)
}

func (u *FeedUsecase) GetByGeckoID(geckoID string) ([]*sql.Feed, error) {
	return u.repo.FeedsByGeckoID(geckoID, 0)
}

func (u *FeedUsecase) GetByDate(date string) ([]*sql.Feed, error) {
	return u.repo.FeedsByDate(date)
}

func (u *FeedUsecase) Delete(geckoID string) error {
	return u.repo.ClearFeed(geckoID)
}
