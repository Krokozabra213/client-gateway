package handler

import (
	"context"
	"fmt"
	"html"
	"strings"

	"github.com/GargantuaLabs/client-gateway/internal/domain"
	authDomain "github.com/GargantuaLabs/client-gateway/internal/features/auth/domain"
	"github.com/GargantuaLabs/client-gateway/internal/features/bot/constants"
	"github.com/GargantuaLabs/client-gateway/internal/features/bot/middlewares"
	"gopkg.in/telebot.v4"
)

func (h *Handler) SubscriptionsCallback(ctx telebot.Context) error {
	if ctx.Callback() != nil {
		_ = ctx.Respond() //nolint:errcheck
	}

	user, ok := ctx.Get(middlewares.CtxUserKey).(*authDomain.User)
	if !ok {
		return ctx.Send("Не удалось определить пользователя. Попробуйте заново.")
	}

	subs, err := h.orderProduct.ListUserSubscriptions(
		context.Background(),
		user.ID.String(),
	)
	if err != nil {
		h.log.Warn("failed to list subscriptions", "error", err.Error())
		return ctx.Send("Не удалось загрузить подписки. Попробуйте позже.")
	}

	if len(subs) == 0 {
		markup := &telebot.ReplyMarkup{}
		shopBtn := markup.Data("🛒 Перейти в каталог", string(constants.ShopCallback))
		backBtn := markup.Data("⬅️ Назад", string(constants.ShowMainMenuCallback))
		markup.Inline(
			markup.Row(shopBtn),
			markup.Row(backBtn),
		)

		return ctx.Send("⭐ <b>У вас пока нет подписок</b>\n\nПерейдите в каталог, чтобы оформить подписку!", &telebot.SendOptions{
			ParseMode:   telebot.ModeHTML,
			ReplyMarkup: markup,
		})
	}

	var sb strings.Builder
	sb.WriteString("⭐ <b>Ваши подписки</b>\n")

	for _, s := range subs {
		sb.WriteString("\n──────────────\n\n")
		sb.WriteString(renderSubscriptionCard(s))
	}

	markup := &telebot.ReplyMarkup{}
	backBtn := markup.Data("⬅️ В меню", string(constants.ShowMainMenuCallback))
	markup.Inline(markup.Row(backBtn))

	return ctx.Send(sb.String(), &telebot.SendOptions{
		ParseMode:   telebot.ModeHTML,
		ReplyMarkup: markup,
	})
}

func renderSubscriptionCard(s domain.Subscription) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%s <b>%s</b>\n",
		subscriptionStatusIcon(s.Status),
		html.EscapeString(s.ProductName),
	))
	sb.WriteString(fmt.Sprintf("Статус: %s\n", subscriptionStatusLabel(s.Status)))
	sb.WriteString(fmt.Sprintf("Истекает: %s UTC\n", s.ExpiresAt.Format("02.01.2006 15:04")))
	sb.WriteString(fmt.Sprintf("ID: <code>%s</code>\n", s.ID))

	return sb.String()
}

func subscriptionStatusIcon(s domain.SubscriptionStatus) string {
	switch s {
	case domain.SubscriptionStatusActive:
		return "🟢"
	case domain.SubscriptionStatusExpired:
		return "🔴"
	case domain.SubscriptionStatusCanceled:
		return "⚫"
	default:
		return "❓"
	}
}

func subscriptionStatusLabel(s domain.SubscriptionStatus) string {
	switch s {
	case domain.SubscriptionStatusActive:
		return "🟢 Активна"
	case domain.SubscriptionStatusExpired:
		return "🔴 Истекла"
	case domain.SubscriptionStatusCanceled:
		return "⚫ Отменена"
	default:
		return "❓ Неизвестно"
	}
}
