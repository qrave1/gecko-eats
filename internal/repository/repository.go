package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/qrave1/gecko-eats/internal/infrastructure/sql"
)

type Repository interface {
	// Pets section
	Pets() ([]sql.Pet, error)
	PetByID(id string) (*sql.Pet, error)
	PetByName(name string) (*sql.Pet, error)
	AddPet(pet *sql.Pet) error

	// Feedings section
	Feedings(petID string, limit int) ([]sql.Feeding, error)
	FeedingsByDate(date string) ([]sql.Feeding, error)
	AddFeeding(petID, date, foodType string) error

	// Others section
	IsWhitelisted(id int64) bool
}

type SqlxRepository struct {
	db        *sqlx.DB
	whitelist []int64
}

var _ Repository = (*SqlxRepository)(nil)

func NewSqlxRepository(db *sqlx.DB, whitelist []int64) *SqlxRepository {
	return &SqlxRepository{db: db, whitelist: whitelist}
}

func (r *SqlxRepository) Pets() ([]sql.Pet, error) {
	var pets []sql.Pet

	err := r.db.Select(&pets, "SELECT * FROM pets ORDER BY name")

	return pets, err
}

func (r *SqlxRepository) PetByName(name string) (*sql.Pet, error) {
	var pet sql.Pet

	err := r.db.Get(&pet, "SELECT * FROM pets WHERE name = ?", name)

	if err != nil {
		return nil, err
	}

	return &pet, nil
}

func (r *SqlxRepository) PetByID(ID string) (*sql.Pet, error) {
	var pet sql.Pet

	err := r.db.Get(&pet, "SELECT * FROM pets WHERE id = ?", ID)

	if err != nil {
		return nil, err
	}

	return &pet, nil
}

func (r *SqlxRepository) AddPet(pet *sql.Pet) error {
	_, err := r.db.Exec("INSERT INTO pets (id, name, food_cycle) VALUES (?, ?, ?)", pet.ID, pet.Name, pet.FoodCycle)

	if err != nil {
		return err
	}

	return nil
}

func (r *SqlxRepository) Feedings(petID string, limit int) ([]sql.Feeding, error) {
	var list []sql.Feeding

	err := r.db.Select(
		&list,
		`
			SELECT * FROM feedings
			WHERE pet_id = ?
			ORDER BY date ASC
			LIMIT ?
				`,
		petID,
		limit,
	)

	return list, err
}

func (r *SqlxRepository) FeedingsByDate(date string) ([]sql.Feeding, error) {
	var list []sql.Feeding

	err := r.db.Select(
		&list,
		`
			SELECT * FROM feedings
			WHERE date = ?
			ORDER BY pet_id ASC
				`,
		date,
	)

	return list, err
}

func (r *SqlxRepository) AddFeeding(petID, date, foodType string) error {
	_, err := r.db.Exec("INSERT INTO feedings (pet_id, date, food_type) VALUES (?, ?, ?)", petID, date, foodType)

	if err != nil {
		return err
	}

	return nil
}

func (r *SqlxRepository) IsWhitelisted(id int64) bool {
	for _, userID := range r.whitelist {
		if userID == id {
			return true
		}
	}

	return false
}
