package telegram

import (
	"fmt"
	"log/slog"
	"strings"

	tele "gopkg.in/telebot.v4"
)

func (b *BotServer) geckoListHandler(c tele.Context) error {
	geckos, err := b.geckoUsecase.GetAll()

	if err != nil {
		slog.Error("Ошибка получения питомцев", "error", err)
		return c.Edit("Ошибка получения питомцев")
	}

	if len(geckos) == 0 {
		return c.Edit("Список пуст")
	}

	return c.Edit("Выберите питомца:", createGeckosKeyboard(geckos))
}

func (b *BotServer) newGeckoHandler(c tele.Context) error {
	gecko, err := b.geckoUsecase.Create(strings.Split(c.Message().Text, " ")[0])

	if err != nil {
		slog.Error("Ошибка добавления питомца", "error", err)

		return c.Send("Ошибка добавления питомца", &tele.SendOptions{ReplyTo: c.Message()})
	}

	err = c.Send(fmt.Sprintf("Питомец %s добавлен!", gecko.Name), &tele.SendOptions{ReplyTo: c.Message()})

	if err != nil {
		slog.Error("Ошибка отправки сообщения", "error", err)

		return err
	}

	return nil
}
