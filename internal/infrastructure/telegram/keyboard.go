package telegram

import (
	"github.com/qrave1/gecko-eats/internal/infrastructure/sql"
	tele "gopkg.in/telebot.v4"
)

var (
	menu         = &tele.ReplyMarkup{ResizeKeyboard: true}
	addGeckoMenu = &tele.ReplyMarkup{}

	btnViewGeckos = menu.Data("📋 Питомцы", "viewGeckos")
	btnAddGecko   = menu.Data("➕ Добавить питомца", "addGecko")
	btnAddfeed    = menu.Data("➕ Добавить календарь", "addfeed")
)

// TODO: клавиатура с возвратом на главное меню

func createGeckosKeyboard(geckos []*sql.Gecko) *tele.ReplyMarkup {
	geckoSelector := &tele.ReplyMarkup{ResizeKeyboard: true}

	var rows []tele.Row

	for _, p := range geckos {
		btn := geckoSelector.Data(p.Name, p.ID)

		rows = append(rows, geckoSelector.Row(btn))
	}

	geckoSelector.Inline(rows...)

	return geckoSelector
}
