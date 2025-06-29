package telegram

import (
	"fmt"
	"strings"

	tele "gopkg.in/telebot.v4"
)

var (
	btnViewPets     = tele.Btn{Text: "📋 Питомцы"}
	btnAddPet       = tele.Btn{Text: "➕ Добавить"}
	inlineMainMenu  = &tele.ReplyMarkup{}
	inlinePetSelect = &tele.ReplyMarkup{}
)

func (b *BotServer) registerHandlers() {
	inlineMainMenu.Inline(
		inlineMainMenu.Row(btnViewPets, btnAddPet),
	)

	b.bot.Handle("/start", func(c tele.Context) error {
		return c.Send("Выберите действие:", inlineMainMenu)
	})

	// показать список питомцев
	b.bot.Handle(&btnViewPets, func(c tele.Context) error {
		pets, err := b.repo.Pets()

		if err != nil {
			return c.Send("Ошибка получения питомцев")
		}

		if len(pets) == 0 {
			return c.Send("Список пуст")
		}

		inlinePetSelect = &tele.ReplyMarkup{}

		var rows []tele.Row

		for _, p := range pets {
			btn := inlinePetSelect.Data(p.Name, "petView", p.ID)

			rows = append(rows, inlinePetSelect.Row(btn))
		}

		inlinePetSelect.Inline(rows...)

		return c.Send("Выберите питомца:", inlinePetSelect)
	})

	// при выборе питомца — показать расписание
	b.bot.Handle(inlinePetSelect.Data, func(c tele.Context) error {
		petID := c.Data()

		pet, err := b.repo.PetByID(petID)

		if err != nil {
			return c.Respond(&tele.CallbackResponse{Text: "Ошибка получения питомца"})
		}

		feedings, err := b.repo.Feedings(pet.ID, 10)

		if err != nil {
			return c.Respond(&tele.CallbackResponse{Text: "Ошибка загрузки расписания"})
		}
		if len(feedings) == 0 {
			return c.Send("Нет запланированных кормлений.")
		}

		var answerBuilder strings.Builder

		answerBuilder.WriteString(fmt.Sprintf("📅 Расписание для %s:\n", pet.Name))

		for _, f := range feedings {
			answerBuilder.WriteString(fmt.Sprintf("• %s — %s\n", f.Date, f.FoodType))
		}

		return c.Send(answerBuilder.String())
	})
}
