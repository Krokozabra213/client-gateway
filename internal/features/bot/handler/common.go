package handler

import (
	"context"
	"fmt"
	"html"
	"time"

	d "github.com/GargantuaLabs/client-gateway/internal/features/auth/domain"
	"github.com/GargantuaLabs/client-gateway/internal/features/bot/constants"
	"github.com/GargantuaLabs/client-gateway/internal/features/shared/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gopkg.in/telebot.v4"
)

func (h *Handler) sendAndStoreMessage(ctx telebot.Context, senderID types.TelegramID, data any, ttl *time.Duration, opts ...any) error {
	msg, err := ctx.Bot().Send(ctx.Chat(), data, opts...)
	if err != nil {
		return err
	}

	return h.cache.StoreMessageID(context.Background(), senderID, msg.ID, ttl)
}

func (h *Handler) removeStoredMessage(ctx telebot.Context, senderID types.TelegramID) error {
	msgId, err := h.cache.GetAndDeleteMessageID(context.Background(), senderID)
	if err != nil {
		return err
	}

	if msgId != 0 {
		msg := &telebot.Message{
			ID:   msgId,
			Chat: ctx.Chat(),
		}

		if err := ctx.Bot().Delete(msg); err != nil {
			h.log.Warn("could not remove message from chat",
				"chat", ctx.Chat().Username,
				"message_id", msg.ID,
				"err", err,
			)
		}
	}

	return nil
}

func (h *Handler) buildMainMenu(user *d.User) (string, *telebot.ReplyMarkup, error) {
	displayName := user.TelegramFirstName

	if user.Username != nil {
		displayName = *user.Username
	}

	profileLink := fmt.Sprintf("tg://user?id=%d", user.TelegramID)
	botCount, err := h.service.GetBotCount(context.Background(), user.ID)

	var currentBalance int64
	resp, err := h.balance.GetAccount(context.Background(), user.ID.String())
	if err != nil {
		st, ok := status.FromError(err)
		if !ok || st.Code() != codes.NotFound {
			return "", nil, fmt.Errorf("failed to get balance: %w", err)
		}
	} else {
		currentBalance = resp.GetBalance()
	}

	if err != nil {
		return "", nil, err
	}

	balanceLine := fmt.Sprintf("💰 Баланс: %.2f ₽", float64(currentBalance)/100)
	warningLine := ""
	if currentBalance <= 0 {
		warningLine = "\n❗️ Баланс пуст — сервисы остановлены"
	}

	text := fmt.Sprintf(`
📚 <a href="%s">%s</a>, твой профиль:

%s%s

🤖 Боты: %d
🧩 Модули: 0/22
`, profileLink, displayName, balanceLine, warningLine, botCount)

	var markup *telebot.ReplyMarkup
	if botCount > 0 {
		markup = h.markups[constants.MainMenuWithBotsMarkup]
	} else {
		markup = h.markups[constants.MainMenuMarkup]
	}

	return text, markup, nil
}

func (h *Handler) mainMenuCallback(ctx telebot.Context) error {
	if ctx.Callback() != nil {
		defer ctx.Respond() //nolint:errcheck
	}

	telegramID, err := types.NewTelegramID(ctx.Sender().ID)
	if err != nil {
		return nil
	}

	user, err := h.authService.GetUserByTelegramID(context.Background(), telegramID)
	if err != nil {
		return nil
	}

	message, markup, err := h.buildMainMenu(user)
	if err != nil {
		return err
	}

	return ctx.Send(message, markup)
}

func formatPrice(amount int64, currency string) string {
	return fmt.Sprintf("%.2f %s", float64(amount)/100, html.EscapeString(currency))
}

func formatDuration(days int) string {
	switch {
	case days <= 0:
		return "—"
	case days%30 == 0:
		return fmt.Sprintf("%d мес.", days/30)
	case days == 1:
		return "1 день"
	case days >= 2 && days <= 4:
		return fmt.Sprintf("%d дня", days)
	default:
		return fmt.Sprintf("%d дней", days)
	}
}
