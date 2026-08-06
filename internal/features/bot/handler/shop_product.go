package handler

import (
	"context"
	"errors"
	"fmt"
	"html"
	"strings"

	"github.com/GargantuaLabs/client-gateway/internal/domain"
	"github.com/GargantuaLabs/client-gateway/internal/features/bot/constants"
	"gopkg.in/telebot.v4"
)

func (h *Handler) ShopProductCallback(ctx telebot.Context) error {
	if ctx.Callback() != nil {
		_ = ctx.Respond() //nolint:errcheck
	}

	productID := ctx.Callback().Data
	if productID == "" {
		return ctx.Send("Товар не найден.")
	}

	product, err := h.orderProduct.GetProduct(context.Background(), productID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return ctx.Send("Товар не найден")
		}
	}

	markup := &telebot.ReplyMarkup{}
	var rows []telebot.Row

	if product.Type == domain.ProductTypeBundle && len(product.BundleItems) > 0 {
		for _, item := range product.BundleItems {
			btn := markup.Data(
				fmt.Sprintf("📋 %s", item.Name),
				string(constants.ShopItemCallback),
				item.ProductID,
			)
			rows = append(rows, markup.Row(btn))
		}
	}

	buyBtn := markup.Data(
		"💳 Купить",
		string(constants.BuyProductCallback),
		product.ID,
	)
	rows = append(rows, markup.Row(buyBtn))

	backBtn := markup.Data(
		"⬅️ К каталогу",
		string(constants.ShopCallback),
	)
	rows = append(rows, markup.Row(backBtn))

	markup.Inline(rows...)

	return ctx.Send(renderProductFull(*product), &telebot.SendOptions{
		ParseMode:   telebot.ModeHTML,
		ReplyMarkup: markup,
	})
}

func renderProductFull(p domain.Product) string {
	var sb strings.Builder

	icon := "📦"
	typeLabel := "Товар"

	switch p.Type {
	case domain.ProductTypeSubscription:
		icon = "⭐"
		typeLabel = "Подписка"
	case domain.ProductTypeBundle:
		icon = "🎁"
		typeLabel = "Пакет"
	}

	sb.WriteString(fmt.Sprintf("%s <b>%s</b>\n\n", icon, html.EscapeString(p.Name)))
	sb.WriteString(fmt.Sprintf("Тип: <b>%s</b>\n", typeLabel))

	if p.Description != "" {
		sb.WriteString(fmt.Sprintf("Описание: %s\n", html.EscapeString(p.Description)))
	}

	if p.DurationDays > 0 {
		sb.WriteString(fmt.Sprintf("Срок: <b>%s</b>\n", formatDuration(p.DurationDays)))
	}

	sb.WriteString(fmt.Sprintf("Цена: <b>%s</b>\n", formatPrice(p.PriceAmount, p.PriceCurrency)))

	if p.Type == domain.ProductTypeBundle && len(p.BundleItems) > 0 {
		sb.WriteString("\n<b>В пакет входит:</b>\n")
		for _, item := range p.BundleItems {
			sb.WriteString(fmt.Sprintf("• %s\n", html.EscapeString(item.Name)))
		}
	}

	return sb.String()
}
