package sql

import "github.com/google/uuid"

type Pet struct {
	ID        string `db:"id"`
	Name      string `db:"name"`
	FoodCycle string `db:"food_cycle"` // JSON-строка: ["обычный", "обычный", "спец"]
}

func NewPet(name string) *Pet {
	return &Pet{ID: uuid.New().String(), Name: name, FoodCycle: `["Кальций", "Кальций", "Кальций","Витамины"]`}
}

type Feeding struct {
	ID       int    `db:"id"`
	Date     string `db:"date"`      // формат YYYY-MM-DD
	PetID    string `db:"pet_id"`    // имя питомца
	FoodType string `db:"food_type"` // строковое описание типа пищи
}
