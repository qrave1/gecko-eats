package telegram

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	tele "gopkg.in/telebot.v4"
)

func (b *BotServer) registerHandlers() {
	menu.Inline(
		menu.Row(btnViewGeckos, btnAddGecko, btnAddfeed),
	)

	b.bot.Handle(
		"/start", func(c tele.Context) error {
			return c.Send("Главное меню:", menu)
		},
	)

	// показать список питомцев
	b.bot.Handle(&btnViewGeckos, b.geckoListHandler)

	b.bot.Handle(
		&btnAddGecko, func(c tele.Context) error {
			return c.Edit("Введите имя питомца в ответ на это сообщение для добавления питомца", addGeckoMenu)
		},
	)

	b.bot.Handle(
		&btnAddfeed, func(c tele.Context) error {
			return c.Edit("Введите имя питомца и интервал кормления в формате: <имя> <интервал>", addGeckoMenu)
		},
	)

	b.bot.Handle(
		tele.OnText, func(c tele.Context) error {
			if c.Message().IsReply() {
				originalMsg := c.Message().ReplyTo.Text

				if strings.Contains(
					originalMsg,
					"Введите имя питомца в ответ на это сообщение для добавления питомца",
				) {
					return b.newGeckoHandler(c)
				} else if strings.Contains(
					originalMsg,
					"Введите имя питомца и интервал кормления в формате: <имя> <интервал>",
				) {
					return b.newfeed(c)
				}
			} else {
				// если не в ответ на сообщение, то просто игнорируем
				slog.Info("Получено сообщение без ответа", "text", c.Text())

				return c.Send("Пожалуйста, ответьте на сообщение для добавления питомца или расписания")
			}

			return nil
		},
	)

	// при выборе питомца — показать расписание
	b.bot.Handle(
		tele.OnCallback, func(c tele.Context) error {
			geckoID := c.Callback().Data

			geckoID = strings.TrimPrefix(geckoID, "\f")

			gecko, err := b.geckoUsecase.GetByID(geckoID)

			if err != nil {
				slog.Error("Ошибка получения питомца", "error", err, "geckoID", geckoID)

				return c.Respond(&tele.CallbackResponse{Text: "Ошибка получения питомца"})
			}

			feeds, err := b.feedUsecase.GetByGeckoID(gecko.ID)

			if err != nil {
				slog.Error("Ошибка получения расписания", "error", err, "geckoID", geckoID)

				return c.Respond(&tele.CallbackResponse{Text: "Ошибка получения расписания"})
			}

			var answerBuilder strings.Builder

			answerBuilder.WriteString(fmt.Sprintf("📅 Расписание для %s:\n", gecko.Name))

			for _, f := range feeds {
				answerBuilder.WriteString(fmt.Sprintf("• %s — %s\n", f.Date, f.FoodType))
			}

			return c.Edit(answerBuilder.String())
		},
	)
}

func (b *BotServer) newfeed(c tele.Context) error {
	parts := strings.Split(c.Text(), " ")

	if len(parts) != 2 {
		return c.Send(
			"Неверный формат. Введите имя питомца и интервал кормления в формате: <имя> <интервал>",
			&tele.SendOptions{ReplyTo: c.Message()},
		)
	}

	gecko, err := b.geckoUsecase.GetByName(parts[0])

	if err != nil {
		slog.Error(
			"create feed schedule: get gecko by name",
			"error", err,
			"geckoName", parts[0],
			"parts", parts,
		)

		return c.Send("Ошибка при получении питомца", &tele.SendOptions{ReplyTo: c.Message()})
	}

	var foodCycleJSON []string

	err = json.Unmarshal([]byte(gecko.FoodCycle), &foodCycleJSON)

	if err != nil {
		slog.Error(
			"create feed schedule: unmarshal food cycle",
			"error", err,
			"food_cycle", gecko.FoodCycle,
		)

		return c.Send(
			"Ошибка при обработке расписания кормления: "+err.Error(),
			&tele.SendOptions{ReplyTo: c.Message()},
		)
	}

	err = b.feedUsecase.DeleteAll(gecko.ID)

	if err != nil {
		slog.Error("Ошибка очистки кормлений", "error", err, "geckoID", gecko.ID)

		return c.Send("Ошибка при очистке кормлений", &tele.SendOptions{ReplyTo: c.Message()})
	}

	interval, err := strconv.Atoi(parts[1])

	if err != nil {
		slog.Error("Ошибка при преобразовании интервала кормления", "error", err, "interval", parts[1])

		return c.Send(
			"Неверный формат интервала. Введите число в формате: <имя> <интервал>",
			&tele.SendOptions{ReplyTo: c.Message()},
		)
	}

	createdfeeds := b.generatefeedSchedule(time.Now(), gecko.ID, foodCycleJSON, interval)

	return c.Send(
		fmt.Sprintf("Расписание кормления создано: %d кормлений", createdfeeds),
		&tele.SendOptions{ReplyTo: c.Message()},
	)
}

func (b *BotServer) generatefeedSchedule(
	startDate time.Time,
	geckoID string,
	geckoCycle []string,
	interval int,
) (createdfeedsCount int) {
	endDate := startDate.AddDate(0, 1, 0) // +1 месяц

	lastFoodIx := 0

	for d := startDate; d.Before(endDate); d = d.AddDate(0, 0, interval) {
		// TODO: добавить проверку на корректность даты

		// TODO: бред какой-то, мб позже понадобится
		//for i, foodType := range geckoCycle {
		//	parsedFoodType, err := time.Parse("2006-01-02", foodType)
		//
		//	if err != nil {
		//		slog.Error("Ошибка при парсинге даты кормления", "error", err, "foodType", foodType)
		//
		//		continue // Пропускаем некорректные даты
		//	}
		//
		//	if parsedFoodType.Year() == d.Year() &&
		//		parsedFoodType.Month() == d.Month() &&
		//		parsedFoodType.Day() == d.Day() {
		//
		//		lastIx = i
		//		break
		//	}
		//}

		// Вечернее кормление (18:00)
		event := time.Date(d.Year(), d.Month(), d.Day(), 18, 0, 0, 0, d.Location())

		err := b.feedUsecase.Create(geckoID, event.Format("2006-01-02"), geckoCycle[lastFoodIx])

		if err != nil {
			slog.Error(
				"add feed",
				slog.Any("error", err),
				slog.String("gecko_id", geckoID),
				slog.String("date", event.Format("2006-01-02")),
				slog.String("food_type", geckoCycle[lastFoodIx]),
			)
		}

		createdfeedsCount++

		// счётчик для следующего типа еды
		if lastFoodIx < len(geckoCycle)-1 {
			lastFoodIx++
		} else {
			lastFoodIx = 0
		}
	}

	return createdfeedsCount
}
