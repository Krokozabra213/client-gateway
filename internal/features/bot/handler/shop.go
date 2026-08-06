package handler

import (
	"context"
	"fmt"

	"github.com/GargantuaLabs/client-gateway/internal/domain"
	"github.com/GargantuaLabs/client-gateway/internal/features/bot/constants"
	"gopkg.in/telebot.v4"
)

func (h *Handler) ShopCallback(ctx telebot.Context) error {
	if ctx.Callback() != nil {
		_ = ctx.Respond() //nolint:errcheck
	}

	products, err := h.orderProduct.ListProducts(context.Background(), true)
	if err != nil {
		return ctx.Send("Не удалось загрузить каталог товаров. Попробуйте позже.")
	}
	if len(products) == 0 {
		return ctx.Send("В каталоге пока нет товаров.")
	}

	markup := &telebot.ReplyMarkup{}
	var rows []telebot.Row

	for _, p := range products {
		icon := "📦"
		if p.Type == domain.ProductTypeSubscription {
			icon = "⭐"
		} else if p.Type == domain.ProductTypeBundle {
			icon = "🎁"
		}

		btnText := fmt.Sprintf("%s %s — %s (%s)",
			icon,
			p.Name,
			formatPrice(p.PriceAmount, p.PriceCurrency),
			formatDuration(p.DurationDays),
		)

		btn := markup.Data(btnText, string(constants.ShopItemCallback), p.ID)

		rows = append(rows, markup.Row(btn))
	}

	backBtn := markup.Data("⬅️ Назад", string(constants.ShowMainMenuCallback))
	rows = append(rows, markup.Row(backBtn))
	markup.Inline(rows...)

	msg := "🛒 <b>Каталог товаров</b>\n\n<i>Нажмите на товар, чтобы посмотреть подробности и состав:</i>"

	return ctx.Send(msg, &telebot.SendOptions{
		ParseMode:   telebot.ModeHTML,
		ReplyMarkup: markup,
	})
}

//func (h *Handler) ShopCallback(ctx telebot.Context) error {
//	if ctx.Callback() != nil {
//		_ = ctx.Respond() //nolint:errcheck
//	}
//
//	products, err := h.orderProduct.ListProducts(context.Background(), true)
//	if err != nil {
//		h.log.Warn("failed to list products", "error", err.Error())
//		return ctx.Send("Не удалось загрузить каталог товаров. Попробуйте позже.")
//	}
//
//	if len(products) == 0 {
//		return ctx.Send("В каталоге пока нет товаров.")
//	}
//
//	markup := &telebot.ReplyMarkup{}
//	rows := make([]telebot.Row, 0, len(products)+1)
//
//	var sb strings.Builder
//	sb.WriteString("🛒 <b>Каталог товаров</b>\n\n")
//
//	for _, p := range products {
//		sb.WriteString(renderProductShort(p))
//		sb.WriteString("\n")
//
//		btn := markup.Data(
//			"Подробнее",
//			string(constants.ShopItemCallback),
//			p.ID,
//		)
//		rows = append(rows, markup.Row(btn))
//	}
//
//	backBtn := markup.Data("⬅️ Назад", string(constants.ShowMainMenuCallback))
//	rows = append(rows, markup.Row(backBtn))
//	markup.Inline(rows...)
//
//	return ctx.Send(sb.String(), &telebot.SendOptions{
//		ParseMode:   telebot.ModeHTML,
//		ReplyMarkup: markup,
//	})
//}
//
//func renderProductShort(p domain.Product) string {
//	icon := "📦"
//	if p.Type == domain.ProductTypeBundle {
//		icon = "🎁"
//	}
//
//	return fmt.Sprintf(
//		"%s <b>%s</b>\n💰 %s\n⏳ %s",
//		icon,
//		html.EscapeString(p.Name),
//		formatPrice(p.PriceAmount, p.PriceCurrency),
//		formatDuration(p.DurationDays),
//	)
//}
