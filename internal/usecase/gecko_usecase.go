package usecase

import (
	"github.com/google/uuid"
	"github.com/qrave1/gecko-eats/internal/infrastructure/sql"
	"github.com/qrave1/gecko-eats/internal/repository"
)

type GeckoUsecase struct {
	repo repository.Repository
}

func NewGeckoUsecase(repo repository.Repository) *GeckoUsecase {
	return &GeckoUsecase{
		repo: repo,
	}
}

func (u *GeckoUsecase) Create(name string) (*sql.Gecko, error) {
	gecko := &sql.Gecko{
		ID:        uuid.New().String(),
		Name:      name,
		FoodCycle: `["Кальций", "Кальций", "Кальций","Витамины"]`,
	}

	if err := u.repo.AddGecko(gecko); err != nil {
		return nil, err
	}

	return gecko, nil
}

func (u *GeckoUsecase) GetByID(geckoID string) (*sql.Gecko, error) {
	return u.repo.GeckoByID(geckoID)
}

func (u *GeckoUsecase) GetByName(name string) (*sql.Gecko, error) {
	return u.repo.GeckoByName(name)
}

func (u *GeckoUsecase) GetAll() ([]*sql.Gecko, error) {
	return u.repo.Geckos()
}
