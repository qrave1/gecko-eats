package telegram

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/qrave1/gecko-eats/internal/domain/input"
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
			return c.Edit(
				"Введите имя питомца и интервал кормления в формате: <имя> <интервал> <дата начала> <индекс первого типа еды>",
				addGeckoMenu,
			)
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
					"Введите имя питомца и интервал кормления в формате: <имя> <интервал> <дата начала> <индекс первого типа еды>",
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
	var err error
	parts := strings.Split(c.Text(), " ")

	if len(parts) < 2 || len(parts) > 4 {
		return c.Send(
			"Неверный формат. Введите имя питомца и интервал кормления в формате: <имя> <интервал> <дата начала> <индекс первого типа еды>",
			&tele.SendOptions{ReplyTo: c.Message()},
		)
	}

	createFeedingsInput := new(input.CreateFeedingsInput)

	createFeedingsInput.Name = parts[0]

	createFeedingsInput.Interval, err = strconv.Atoi(parts[1])
	if err != nil {
		slog.Error("Ошибка при преобразовании интервала кормления", "error", err, "interval", parts[1])

		return c.Send(
			"Неверный формат интервала. Введите число в формате: <имя> <интервал> <дата начала> <индекс первого типа еды>",
			&tele.SendOptions{ReplyTo: c.Message()},
		)
	}

	createFeedingsInput.StartDate = time.Now()

	switch len(parts) {
	case 2:
	case 3:
		createFeedingsInput.StartDate, err = time.Parse("02-01-2006", parts[2])
		if err != nil {
			slog.Error("Ошибка при преобразовании даты начала", "error", err)

			return c.Send(
				"Дата написана неправильно",
				&tele.SendOptions{ReplyTo: c.Message()},
			)
		}
	case 4:
		createFeedingsInput.StartDate, err = time.Parse("02-01-2006", parts[2])
		if err != nil {
			slog.Error("Ошибка при преобразовании даты начала", "error", err)

			return c.Send(
				"Дата написана неправильно",
				&tele.SendOptions{ReplyTo: c.Message()},
			)
		}

		createFeedingsInput.StartIx, err = strconv.Atoi(parts[3])
	default:
		return c.Send(
			"Неправильное количество параметров",
			&tele.SendOptions{ReplyTo: c.Message()},
		)
	}

	gecko, err := b.geckoUsecase.GetByName(createFeedingsInput.Name)
	if err != nil {
		slog.Error(
			"create feed schedule: get gecko by name",
			"error", err,
			"geckoName", createFeedingsInput.Name,
			"parts", parts,
		)

		return c.Send("Ошибка при получении питомца", &tele.SendOptions{ReplyTo: c.Message()})
	}

	err = b.feedUsecase.DeleteAll(gecko.ID)
	if err != nil {
		slog.Error("Ошибка очистки кормлений", "error", err, "geckoID", gecko.ID)

		return c.Send("Ошибка при очистке кормлений", &tele.SendOptions{ReplyTo: c.Message()})
	}

	createdFeeds, err := b.feedUsecase.CreateSchedule(gecko, createFeedingsInput)
	if err != nil {
		slog.Error("Ошибка создания расписания кормления", "error", err, "interval", parts[1])

		return c.Send(
			"Ошибка создания расписания кормления",
			&tele.SendOptions{ReplyTo: c.Message()},
		)
	}

	return c.Send(
		fmt.Sprintf("Расписание кормления создано: %d кормлений", createdFeeds),
		&tele.SendOptions{ReplyTo: c.Message()},
	)
}
