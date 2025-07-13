package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/qrave1/gecko-eats/internal/domain"
)

type Repository interface {
	// Geckos section
	Geckos() ([]*domain.Gecko, error)
	GeckoByID(id string) (*domain.Gecko, error)
	GeckoByName(name string) (*domain.Gecko, error)
	AddGecko(gecko *domain.Gecko) error

	// feeds section
	FeedsByGeckoID(geckoID string, limit int) ([]*domain.Feed, error)
	FeedsByDate(date string) ([]*domain.Feed, error)
	AddFeed(feed *domain.Feed) error
	ClearFeed(geckoID string) error
}

type SqlxRepository struct {
	db *sqlx.DB
}

var _ Repository = (*SqlxRepository)(nil)

func NewSqlxRepository(db *sqlx.DB) *SqlxRepository {
	return &SqlxRepository{db: db}
}

func (r *SqlxRepository) Geckos() ([]*domain.Gecko, error) {
	var geckos []*domain.Gecko

	err := r.db.Select(&geckos, "SELECT * FROM geckos ORDER BY name")

	return geckos, err
}

func (r *SqlxRepository) GeckoByName(name string) (*domain.Gecko, error) {
	var gecko domain.Gecko

	err := r.db.Get(&gecko, "SELECT * FROM geckos WHERE name = $1", name)

	if err != nil {
		return nil, err
	}

	return &gecko, nil
}

func (r *SqlxRepository) GeckoByID(ID string) (*domain.Gecko, error) {
	var gecko domain.Gecko

	err := r.db.Get(&gecko, "SELECT * FROM geckos WHERE id = $1", ID)

	if err != nil {
		return nil, err
	}

	return &gecko, nil
}

func (r *SqlxRepository) AddGecko(gecko *domain.Gecko) error {
	_, err := r.db.Exec("INSERT INTO geckos (id, name, food_cycle) VALUES ($1, $2, $3)", gecko.ID, gecko.Name, gecko.FoodCycle)

	if err != nil {
		return err
	}

	return nil
}

func (r *SqlxRepository) FeedsByGeckoID(geckoID string, limit int) ([]*domain.Feed, error) {
	var list []*domain.Feed

	err := r.db.Select(
		&list,
		`
			SELECT * FROM feeds
			WHERE gecko_id = $1
			ORDER BY date ASC
			LIMIT $2
				`,
		geckoID,
		limit,
	)

	return list, err
}

func (r *SqlxRepository) FeedsByDate(date string) ([]*domain.Feed, error) {
	var list []*domain.Feed

	err := r.db.Select(
		&list,
		`
			SELECT * FROM feeds
			WHERE date = $1
			ORDER BY gecko_id ASC
				`,
		date,
	)

	return list, err
}

func (r *SqlxRepository) AddFeed(feed *domain.Feed) error {
	_, err := r.db.Exec(
		"INSERT INTO feeds (gecko_id, date, food_type) VALUES ($1, $2, $3)",
		feed.GeckoID,
		feed.Date,
		feed.FoodType,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *SqlxRepository) ClearFeed(geckoID string) error {
	_, err := r.db.Exec("DELETE FROM feeds WHERE gecko_id = $1", geckoID)

	if err != nil {
		return err
	}

	return nil
}
