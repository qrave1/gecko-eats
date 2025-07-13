package usecase

import (
	"github.com/google/uuid"
	"github.com/qrave1/gecko-eats/internal/domain"
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

func (u *GeckoUsecase) Create(name string) (*domain.Gecko, error) {
	gecko := &domain.Gecko{
		ID:        uuid.New().String(),
		Name:      name,
		FoodCycle: `["Кальций", "Кальций", "Кальций", "Витамины"]`,
	}

	if err := u.repo.AddGecko(gecko); err != nil {
		return nil, err
	}

	return gecko, nil
}

func (u *GeckoUsecase) GetByID(geckoID string) (*domain.Gecko, error) {
	return u.repo.GeckoByID(geckoID)
}

func (u *GeckoUsecase) GetByName(name string) (*domain.Gecko, error) {
	return u.repo.GeckoByName(name)
}

func (u *GeckoUsecase) GetAll() ([]*domain.Gecko, error) {
	return u.repo.Geckos()
}
