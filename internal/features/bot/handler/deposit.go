package handler

import (
	"context"
	"fmt"
	"strconv"

	"github.com/GargantuaLabs/client-gateway/internal/domain"
	authDomain "github.com/GargantuaLabs/client-gateway/internal/features/auth/domain"
	"github.com/GargantuaLabs/client-gateway/internal/features/bot/constants"
	"github.com/GargantuaLabs/client-gateway/internal/features/bot/middlewares"
	"gopkg.in/telebot.v4"
)

var starOptions = []int{1, 50, 100, 250, 500, 1000}
var tbankOptions = []int{1, 100, 500, 1000, 2000, 5000}

func (h *Handler) depositCallback(ctx telebot.Context) error {
	if ctx.Callback() != nil {
		defer ctx.Respond() //nolint:errcheck
	}

	return ctx.Send("💰 Выберите способ пополнения", h.markups[constants.DepositMarkup])
}

func (h *Handler) TGDepositCallback(ctx telebot.Context) error {
	if ctx.Callback() != nil {
		_ = ctx.Respond() //nolint:errcheck
	}

	markup := &telebot.ReplyMarkup{}
	var rows []telebot.Row

	// Кнопки с суммами (по 2 в ряд)
	for i := 0; i < len(starOptions); i += 2 {
		var rowBtns []telebot.Btn

		btn1 := markup.Data(
			fmt.Sprintf("⭐ %d", starOptions[i]),
			string(constants.TGDepositAmountCallback),
			strconv.Itoa(starOptions[i]),
		)
		rowBtns = append(rowBtns, btn1)

		if i+1 < len(starOptions) {
			btn2 := markup.Data(
				fmt.Sprintf("⭐ %d", starOptions[i+1]),
				string(constants.TGDepositAmountCallback),
				strconv.Itoa(starOptions[i+1]),
			)
			rowBtns = append(rowBtns, btn2)
		}

		rows = append(rows, markup.Row(rowBtns...))
	}

	// Кнопка назад
	backBtn := markup.Data("⬅️ Назад", string(constants.DepositCallback))
	rows = append(rows, markup.Row(backBtn))

	markup.Inline(rows...)

	msg := `⭐ <b>Пополнение звёздами Telegram</b>

Выберите количество звёзд для пополнения:`

	return ctx.Send(msg, &telebot.SendOptions{
		ParseMode:   telebot.ModeHTML,
		ReplyMarkup: markup,
	})
}

func (h *Handler) tbankDepositCallback(ctx telebot.Context) error {
	if ctx.Callback() != nil {
		_ = ctx.Respond() //nolint:errcheck
	}

	markup := &telebot.ReplyMarkup{}
	var rows []telebot.Row

	for i := 0; i < len(tbankOptions); i += 2 {
		var rowBtns []telebot.Btn

		btn1 := markup.Data(
			fmt.Sprintf("💳 %d ₽", tbankOptions[i]),
			string(constants.TbankDepositAmountCallback),
			strconv.Itoa(tbankOptions[i]),
		)
		rowBtns = append(rowBtns, btn1)

		if i+1 < len(tbankOptions) {
			btn2 := markup.Data(
				fmt.Sprintf("💳 %d ₽", tbankOptions[i+1]),
				string(constants.TbankDepositAmountCallback),
				strconv.Itoa(tbankOptions[i+1]),
			)
			rowBtns = append(rowBtns, btn2)
		}

		rows = append(rows, markup.Row(rowBtns...))
	}

	backBtn := markup.Data("⬅️ Назад", string(constants.DepositCallback))
	rows = append(rows, markup.Row(backBtn))

	markup.Inline(rows...)

	msg := `💳 <b>Пополнение через TBank</b>

Выберите сумму для пополнения:`

	return ctx.Send(msg, &telebot.SendOptions{
		ParseMode:   telebot.ModeHTML,
		ReplyMarkup: markup,
	})
}

func (h *Handler) TGDepositAmountCallback(ctx telebot.Context) error {
	if ctx.Callback() != nil {
		_ = ctx.Respond() //nolint:errcheck
	}

	user, ok := ctx.Get(middlewares.CtxUserKey).(*authDomain.User)
	if !ok {
		return ctx.Send("Не удалось определить пользователя.")
	}

	starsStr := ctx.Data()
	stars, err := strconv.Atoi(starsStr)
	if err != nil || stars <= 0 {
		return ctx.Send("Некорректное количество звёзд.")
	}

	msgID := ctx.Callback().Message.ID
	idempotencyKey := CreateIdempotency(user.ID.String(), msgID)
	tgPayload := domain.NewTgInit(idempotencyKey, user.ID.String(), "Пополнение баланса", "XTR", "Пополнение звёздами", stars)

	resp, err := h.payment.TGInit(context.Background(), tgPayload)
	if err != nil {
		h.log.Warn("failed to create tg payment", "error", err.Error())

		markup := &telebot.ReplyMarkup{}
		backBtn := markup.Data("⬅️ Назад", string(constants.TGDepositCallback))
		markup.Inline(markup.Row(backBtn))

		return ctx.Send("❌ Ошибка при создании платежа. Попробуйте позже.", &telebot.SendOptions{
			ReplyMarkup: markup,
		})
	}

	msg := fmt.Sprintf(`✅ <b>Счёт создан!</b>

Сумма: <b>%d ⭐</b>

Перейдите по ссылке для оплаты:
%s`, stars, resp.PaymentUrl)

	markup := &telebot.ReplyMarkup{}
	backBtn := markup.Data("⬅️ В меню", string(constants.ShowMainMenuCallback))
	markup.Inline(markup.Row(backBtn))

	return ctx.Send(msg, &telebot.SendOptions{
		ParseMode:   telebot.ModeHTML,
		ReplyMarkup: markup,
	})
}

//func (h *Handler) depositTbank(ctx telebot.Context) error {
//	user, ok := ctx.Get(middlewares.CtxUserKey).(*authDomain.User)
//	if !ok {
//		return nil
//	}
//
//	rubs, err := strconv.Atoi(ctx.Text())
//	if err != nil {
//		return ctx.Send("Невалидное значение, попробуйте ещё раз.")
//	}
//
//	msgID := ctx.Message().ID
//	idempotencyKey := CreateIdempotency(user.ID.String(), msgID)
//
//	tbankPayload := domain.NewTbankInit(idempotencyKey, user.ID.String(), "desc", rubs)
//	resp, err := h.payment.TbankInit(context.Background(), tbankPayload)
//	if err != nil {
//		return ctx.Send("Ошибка при создании платежа. Попробуйте позже.")
//	}
//
//	h.cache.RemoveLastAction(context.Background(), user.TelegramID) //nolint:errcheck
//	return ctx.Send(resp.PaymentUrl)
//}

func (h *Handler) TbankDepositAmountCallback(ctx telebot.Context) error {
	if ctx.Callback() != nil {
		_ = ctx.Respond() //nolint:errcheck
	}

	user, ok := ctx.Get(middlewares.CtxUserKey).(*authDomain.User)
	if !ok {
		return ctx.Send("Не удалось определить пользователя.")
	}

	rubsStr := ctx.Data()
	rubs, err := strconv.Atoi(rubsStr)
	if err != nil || rubs <= 0 {
		return ctx.Send("Некорректная сумма.")
	}

	msgID := ctx.Callback().Message.ID
	idempotencyKey := CreateIdempotency(user.ID.String(), msgID)
	tbankPayload := domain.NewTbankInit(idempotencyKey, user.ID.String(), "Пополнение баланса", rubs)

	resp, err := h.payment.TbankInit(context.Background(), tbankPayload)
	if err != nil {
		h.log.Warn("failed to create tbank payment", "error", err.Error())

		markup := &telebot.ReplyMarkup{}
		backBtn := markup.Data("⬅️ Назад", string(constants.TbankDepositCallback))
		markup.Inline(markup.Row(backBtn))

		return ctx.Send("❌ Ошибка при создании платежа. Попробуйте позже.", &telebot.SendOptions{
			ReplyMarkup: markup,
		})
	}

	msg := fmt.Sprintf(`✅ <b>Счёт создан!</b>

Сумма: <b>%d ₽</b>

Перейдите по ссылке для оплаты:
%s`, rubs, resp.PaymentUrl)

	markup := &telebot.ReplyMarkup{}
	backBtn := markup.Data("⬅️ В меню", string(constants.ShowMainMenuCallback))
	markup.Inline(markup.Row(backBtn))

	return ctx.Send(msg, &telebot.SendOptions{
		ParseMode:   telebot.ModeHTML,
		ReplyMarkup: markup,
	})
}

func (h *Handler) HandleTGCheckout(ctx telebot.Context) error {
	checkout := ctx.PreCheckoutQuery()

	checkoutPayload := &domain.PreCheckoutPayload{
		ID:             checkout.ID,
		Currency:       checkout.Currency,
		TotalAmount:    checkout.Total,
		InvoicePayload: checkout.Payload,
		From: domain.User{
			ID:        checkout.Sender.ID,
			IsBot:     checkout.Sender.IsBot,
			FirstName: checkout.Sender.FirstName,
		},
	}

	err := h.payment.TGPreCheckout(context.Background(), checkoutPayload)
	if err != nil {
		return ctx.Send("Ошибка при создании платежа. Попробуйте позже.")
	}
	return nil
}

func (h *Handler) HandleTGSuccessful(ctx telebot.Context) error {
	payment := ctx.Message().Payment

	successfulPayload := &domain.SuccessfulPayload{
		Currency:                payment.Currency,
		InvoicePayload:          payment.Payload,
		TelegramPaymentChargeID: payment.TelegramChargeID,
		ProviderPaymentChargeID: payment.ProviderChargeID,
		TotalAmount:             payment.Total,
	}

	err := h.payment.TGSuccessful(context.Background(), successfulPayload)
	if err != nil {
		return ctx.Send("Ошибка при создании платежа. Попробуйте позже.")
	}
	return nil
}

func CreateIdempotency(userID string, msgID int) string {
	return fmt.Sprintf("dep_%s_%d", userID, msgID)
}
