package handler

import (
	"github.com/GargantuaLabs/client-gateway/internal/features/bot/constants"
	"gopkg.in/telebot.v4"
)

func (h *Handler) helpCallback(ctx telebot.Context) error {
	if ctx.Callback() != nil {
		_ = ctx.Respond() //nolint:errcheck
	}

	msg := `🛠 <b>Раздел «Помощь»</b>

Этот раздел находится в активной разработке ✨

Совсем скоро здесь появятся:
• 📖 Инструкции по настройке ботов
• 💡 Ответы на частые вопросы
• 🎓 Обучающие материалы
• 📩 Форма обратной связи

А пока вы можете написать в поддержку:
👉 @your_support_username

Спасибо за терпение! 🙏`

	markup := &telebot.ReplyMarkup{}
	backBtn := markup.Data("⬅️ Назад", string(constants.ShowMainMenuCallback))
	markup.Inline(markup.Row(backBtn))

	return ctx.Send(msg, &telebot.SendOptions{
		ParseMode:   telebot.ModeHTML,
		ReplyMarkup: markup,
	})
}
