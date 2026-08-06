package handler

import (
	"context"
	"fmt"
	"html"
	"strconv"
	"strings"

	"github.com/GargantuaLabs/client-gateway/internal/domain"
	authDomain "github.com/GargantuaLabs/client-gateway/internal/features/auth/domain"
	"github.com/GargantuaLabs/client-gateway/internal/features/bot/constants"
	"github.com/GargantuaLabs/client-gateway/internal/features/bot/middlewares"
	"gopkg.in/telebot.v4"
)

const ordersPerPage = 5

func (h *Handler) OrdersCallback(ctx telebot.Context) error {
	if ctx.Callback() != nil {
		_ = ctx.Respond() //nolint:errcheck
	}

	user, ok := ctx.Get(middlewares.CtxUserKey).(*authDomain.User)
	if !ok {
		return ctx.Send("Не удалось определить пользователя. Попробуйте заново.")
	}

	offset := 0
	if ctx.Callback() != nil && ctx.Callback().Data != "" {
		if parsed, err := strconv.Atoi(ctx.Callback().Data); err == nil {
			offset = parsed
		}
	}

	orders, err := h.orderProduct.ListUserOrders(
		context.Background(),
		user.ID.String(),
		ordersPerPage+1,
		int32(offset),
	)
	if err != nil {
		h.log.Warn("failed to list orders", "error", err.Error())
		return ctx.Send("Не удалось загрузить заказы. Попробуйте позже.")
	}

	if len(orders) == 0 && offset == 0 {
		markup := &telebot.ReplyMarkup{}
		shopBtn := markup.Data("🛒 Перейти в каталог", string(constants.ShopCallback))
		backBtn := markup.Data("⬅️ Назад", string(constants.ShowMainMenuCallback))
		markup.Inline(
			markup.Row(shopBtn),
			markup.Row(backBtn),
		)

		return ctx.Send("📦 <b>У вас пока нет заказов</b>\n\nПерейдите в каталог, чтобы сделать первый заказ!", &telebot.SendOptions{
			ParseMode:   telebot.ModeHTML,
			ReplyMarkup: markup,
		})
	}

	hasNext := len(orders) > ordersPerPage
	if hasNext {
		orders = orders[:ordersPerPage]
	}
	hasPrev := offset > 0

	var sb strings.Builder
	page := (offset / ordersPerPage) + 1
	sb.WriteString(fmt.Sprintf("📦 <b>Ваши заказы</b> (стр. %d)\n", page))

	for _, o := range orders {
		sb.WriteString("\n──────────────\n\n")
		sb.WriteString(renderOrderCard(o))
	}

	markup := &telebot.ReplyMarkup{}
	var rows []telebot.Row

	var navBtns []telebot.Btn
	if hasPrev {
		prevOffset := offset - ordersPerPage
		if prevOffset < 0 {
			prevOffset = 0
		}
		prevBtn := markup.Data("◀️ Назад", string(constants.MyOrdersCallback), strconv.Itoa(prevOffset))
		navBtns = append(navBtns, prevBtn)
	}
	if hasNext {
		nextOffset := offset + ordersPerPage
		nextBtn := markup.Data("Вперёд ▶️", string(constants.MyOrdersCallback), strconv.Itoa(nextOffset))
		navBtns = append(navBtns, nextBtn)
	}
	if len(navBtns) > 0 {
		rows = append(rows, markup.Row(navBtns...))
	}

	backBtn := markup.Data("⬅️ В меню", string(constants.ShowMainMenuCallback))
	rows = append(rows, markup.Row(backBtn))

	markup.Inline(rows...)

	if ctx.Callback() != nil {
		return ctx.Edit(sb.String(), &telebot.SendOptions{
			ParseMode:   telebot.ModeHTML,
			ReplyMarkup: markup,
		})
	}

	return ctx.Send(sb.String(), &telebot.SendOptions{
		ParseMode:   telebot.ModeHTML,
		ReplyMarkup: markup,
	})
}

func renderOrderCard(o domain.Order) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%s <b>%s</b>\n",
		orderStatusIcon(o.Status),
		html.EscapeString(o.ProductName),
	))
	sb.WriteString(fmt.Sprintf("Тип: %s\n", productTypeLabel(o.ProductType)))
	sb.WriteString(fmt.Sprintf("Срок: <b>%s</b>\n", formatDuration(o.DurationDays)))
	sb.WriteString(fmt.Sprintf("Сумма: <b>%.2f %s</b>\n",
		float64(o.Amount)/100,
		html.EscapeString(o.Currency),
	))
	sb.WriteString(fmt.Sprintf("Статус: %s\n", orderStatusLabel(o.Status)))
	sb.WriteString(fmt.Sprintf("Создан: %s UTC\n", o.CreatedAt.Format("02.01.2006 15:04")))
	sb.WriteString(fmt.Sprintf("ID: <code>%s</code>\n", o.ID))

	return sb.String()
}

func orderStatusIcon(s domain.OrderStatus) string {
	switch s {
	case domain.OrderStatusCompleted:
		return "✅"
	case domain.OrderStatusPending:
		return "⏳"
	case domain.OrderStatusFailed:
		return "❌"
	case domain.OrderStatusCancelled:
		return "🚫"
	case domain.OrderStatusRefunded:
		return "↩️"
	case domain.OrderStatusExpired:
		return "⌛"
	default:
		return "❓"
	}
}

func orderStatusLabel(s domain.OrderStatus) string {
	switch s {
	case domain.OrderStatusCompleted:
		return "✅ Выполнен"
	case domain.OrderStatusPending:
		return "⏳ В обработке"
	case domain.OrderStatusFailed:
		return "❌ Ошибка"
	case domain.OrderStatusCancelled:
		return "🚫 Отменён"
	case domain.OrderStatusRefunded:
		return "↩️ Возврат"
	case domain.OrderStatusExpired:
		return "⌛ Истёк"
	default:
		return "❓ Неизвестен"
	}
}

func productTypeLabel(t domain.ProductType) string {
	switch t {
	case domain.ProductTypeSubscription:
		return "⭐ Подписка"
	case domain.ProductTypeBundle:
		return "🎁 Пакет"
	default:
		return "📦 Товар"
	}
}
