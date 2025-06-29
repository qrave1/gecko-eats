package middleware

import tele "gopkg.in/telebot.v4"

// WhitelistOnly проверяет, находится ли ID пользователя в белом списке
func WhitelistOnly(allowedUsers []int64) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			for _, id := range allowedUsers {
				if id == c.Sender().ID {
					return next(c)
				}
			}

			return c.Send("Извините, у вас нет доступа к этому боту.")
		}
	}
}
