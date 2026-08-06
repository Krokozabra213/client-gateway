package handler

import (
	"github.com/GargantuaLabs/client-gateway/internal/features/bot/constants"
	"gopkg.in/telebot.v4"
)

func (h *Handler) statisticsCallback(ctx telebot.Context) error {
	if ctx.Callback() != nil {
		_ = ctx.Respond() //nolint:errcheck
	}

	msg := `📊 <b>Статистика</b>

Этот раздел находится в активной разработке 📈

Совсем скоро здесь появятся:
• 👥 Аудитория ботов
• 💬 Количество сообщений
• 📉 Графики активности
• 🏆 Рейтинг модулей

Следите за обновлениями! ✨`

	markup := &telebot.ReplyMarkup{}
	backBtn := markup.Data("⬅️ Назад", string(constants.ShowMainMenuCallback))
	markup.Inline(markup.Row(backBtn))

	return ctx.Send(msg, &telebot.SendOptions{
		ParseMode:   telebot.ModeHTML,
		ReplyMarkup: markup,
	})
}
