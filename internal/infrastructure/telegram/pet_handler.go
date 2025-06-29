package telegram

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/qrave1/gecko-eats/internal/infrastructure/sql"
	tele "gopkg.in/telebot.v4"
)

var (
	menu        = &tele.ReplyMarkup{ResizeKeyboard: true}
	petSelector = &tele.ReplyMarkup{}
	addPetMenu  = &tele.ReplyMarkup{}

	btnViewPets   = menu.Data("📋 Питомцы", "viewPets")
	btnAddPet     = menu.Data("➕ Добавить питомца", "addPet")
	btnAddFeeding = menu.Data("➕ Добавить календарь", "addFeeding")
)

func (b *BotServer) registerHandlers() {
	menu.Inline(
		menu.Row(btnViewPets, btnAddPet, btnAddFeeding),
	)

	b.bot.Handle("/start", func(c tele.Context) error {
		return c.Send("Главное меню:", menu)
	})

	// показать список питомцев
	b.bot.Handle(&btnViewPets, func(c tele.Context) error {
		pets, err := b.repo.Pets()

		if err != nil {
			slog.Error("Ошибка получения питомцев", "error", err)

			return c.Edit("Ошибка получения питомцев")
		}

		if len(pets) == 0 {
			return c.Edit("Список пуст")
		}

		petSelector = &tele.ReplyMarkup{}

		var rows []tele.Row

		for _, p := range pets {
			btn := petSelector.Data(p.Name, p.ID)

			rows = append(rows, petSelector.Row(btn))
		}

		petSelector.Inline(rows...)

		return c.Edit("Выберите питомца:", petSelector)
	})

	b.bot.Handle(&btnAddPet, func(c tele.Context) error {
		return c.Edit("Введите имя питомца в ответ на это сообщение для добавления питомца", addPetMenu)
	})

	b.bot.Handle(&btnAddFeeding, func(c tele.Context) error {
		return c.Edit("Введите имя питомца и интервал кормления в формате: <имя> <интервал>", addPetMenu)
	})

	b.bot.Handle(tele.OnText, func(c tele.Context) error {
		if c.Message().IsReply() {
			originalMsg := c.Message().ReplyTo.Text

			if strings.Contains(originalMsg, "Введите имя питомца в ответ на это сообщение для добавления питомца") {
				return b.newPet(c)
			} else if strings.Contains(originalMsg, "Введите имя питомца и интервал кормления в формате: <имя> <интервал>") {
				return b.newFeeding(c)
			}
		}

		return nil
	})

	// при выборе питомца — показать расписание
	b.bot.Handle(tele.OnCallback, func(c tele.Context) error {
		petID := c.Callback().Data

		petID = strings.TrimPrefix(petID, "\f")

		pet, err := b.repo.PetByID(petID)

		if err != nil {
			slog.Error("Ошибка получения питомца", "error", err, "petID", petID)

			return c.Respond(&tele.CallbackResponse{Text: "Ошибка получения питомца"})
		}

		feedings, err := b.repo.Feedings(pet.ID, 10)

		if err != nil {
			slog.Error("Ошибка получения расписания", "error", err, "petID", petID)

			return c.Respond(&tele.CallbackResponse{Text: "Ошибка получения расписания"})
		}

		var answerBuilder strings.Builder

		answerBuilder.WriteString(fmt.Sprintf("📅 Расписание для %s:\n", pet.Name))

		for _, f := range feedings {
			answerBuilder.WriteString(fmt.Sprintf("• %s — %s\n", f.Date, f.FoodType))
		}

		return c.Edit(answerBuilder.String())
	})
}

func (b *BotServer) newPet(c tele.Context) error {
	pet := sql.NewPet(strings.Split(c.Message().Text, " ")[0])

	err := b.repo.AddPet(pet)

	if err != nil {
		slog.Error("Ошибка добавления питомца", "error", err)

		return c.Send("Ошибка добавления питомца", &tele.SendOptions{ReplyTo: c.Message()})
	}

	err = c.Send("Питомец добавлен!", &tele.SendOptions{ReplyTo: c.Message()})

	if err != nil {
		slog.Error("Ошибка отправки сообщения", "error", err)

		return err
	}

	return nil
}

func (b *BotServer) newFeeding(c tele.Context) error {
	parts := strings.Split(c.Text(), " ")

	if len(parts) != 2 {
		return c.Send(
			"Неверный формат. Введите имя питомца и интервал кормления в формате: <имя> <интервал>",
			&tele.SendOptions{ReplyTo: c.Message()},
		)
	}

	pet, err := b.repo.PetByName(parts[0])

	if err != nil {
		slog.Error(
			"create feeding schedule: get pet by name",
			"error", err,
			"petName", parts[0],
			"parts", parts,
		)

		return c.Send("Ошибка при получении питомца", &tele.SendOptions{ReplyTo: c.Message()})
	}

	var foodCycleJSON []string

	err = json.Unmarshal([]byte(pet.FoodCycle), &foodCycleJSON)

	if err != nil {
		slog.Error(
			"create feeding schedule: unmarshal food cycle",
			"error", err,
			"food_cycle", pet.FoodCycle,
		)

		return c.Send("Ошибка при обработке расписания кормления: "+err.Error(), &tele.SendOptions{ReplyTo: c.Message()})
	}

	schedule := b.generateFeedingSchedule(time.Now(), pet.ID, foodCycleJSON)

	return c.Send(
		fmt.Sprintf("Расписание кормления создано: %d кормлений", len(schedule)),
		&tele.SendOptions{ReplyTo: c.Message()},
	)
}

func (b *BotServer) generateFeedingSchedule(
	startDate time.Time,
	petID string,
	petCycle []string,
) []sql.Feeding {
	var schedule []sql.Feeding

	endDate := startDate.AddDate(0, 1, 0) // +1 месяц

	lastIx := 0

	for d := startDate; d.Before(endDate); d = d.AddDate(0, 0, 1) {
		// Проверяем что дата существует (не 31 февраля и т.п.)
		if d.Day() != startDate.Day() && d.Day() < startDate.Day() {
			continue // Пропускаем несуществующие даты
		}

		for i, foodType := range petCycle {
			parsedFoodType, err := time.Parse("2006-01-02", foodType)

			if err != nil {
				slog.Error("Ошибка при парсинге даты кормления", "error", err, "foodType", foodType)

				continue // Пропускаем некорректные даты
			}

			if parsedFoodType.Year() == d.Year() &&
				parsedFoodType.Month() == d.Month() &&
				parsedFoodType.Day() == d.Day() {

				lastIx = i
				break
			}
		}

		// Вечернее кормление (18:00)
		evening := time.Date(d.Year(), d.Month(), d.Day(), 18, 0, 0, 0, d.Location())

		err := b.repo.AddFeeding(petID, evening.Format("2006-01-02"), petCycle[lastIx])

		if err != nil {
			slog.Error(
				"add feeding",
				slog.Any("error", err),
				slog.String("pet_id", petID),
				slog.String("date", evening.Format("2006-01-02")),
				slog.String("food_type", petCycle[lastIx]),
			)
		}
	}

	return schedule
}
