package domain

type Gecko struct {
	ID        string `db:"id"`
	Name      string `db:"name"`
	FoodCycle string `db:"food_cycle"` // JSON-строка: ["обычный", "обычный", "спец"]
}

type Feed struct {
	Date     string `db:"date"`      // формат YYYY-MM-DD
	GeckoID  string `db:"gecko_id"`  // имя питомца
	FoodType string `db:"food_type"` // строковое описание типа пищи
}
