package handler

import (
	"context"
	"fmt"

	authDomain "github.com/GargantuaLabs/client-gateway/internal/features/auth/domain"
	"github.com/GargantuaLabs/client-gateway/internal/features/bot/constants"
	"github.com/GargantuaLabs/client-gateway/internal/features/bot/middlewares"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/telebot.v4"
)

func (h *Handler) BuyProductCallback(ctx telebot.Context) error {
	if ctx.Callback() != nil {
		_ = ctx.Respond() //nolint:errcheck
	}

	productID := ctx.Callback().Data
	if productID == "" {
		return ctx.Send("Товар не найден.")
	}

	user, ok := ctx.Get(middlewares.CtxUserKey).(*authDomain.User)
	if !ok {
		return ctx.Send("Не удалось определить пользователя. Попробуйте заново.")
	}

	//tgID, err := types.NewTelegramID(ctx.Sender().ID)
	//if err != nil {
	//	return ctx.Send("Невалидный telegramID")
	//}
	//
	//user, err := h.authService.GetUserByTelegramID(context.Background(), tgID)
	//if err != nil {
	//	return ctx.Send("Пользователь не найден. Пожалуйста, зарегистрируйтесь.")
	//}

	orderID, _, err := h.orderProduct.PurchaseProduct(
		context.Background(),
		user.ID.String(),
		productID,
	)

	markup := &telebot.ReplyMarkup{}
	backBtn := markup.Data("⬅️ К каталогу", string(constants.ShopCallback))
	balanceBtn := markup.Data("💰 Пополнить баланс", string(constants.DepositCallback))

	if err != nil {
		return h.handlePurchaseError(ctx, err, markup, backBtn, balanceBtn)
	}

	return h.handlePurchaseSuccess(ctx, orderID, markup, backBtn)
}

func (h *Handler) handlePurchaseError(
	ctx telebot.Context,
	err error,
	markup *telebot.ReplyMarkup,
	backBtn, balanceBtn telebot.Btn,
) error {
	st, ok := status.FromError(err)
	if !ok {
		h.log.Warn("failed to purchase product", "error", err.Error())
		markup.Inline(markup.Row(backBtn))
		return ctx.Send("⚠️ Не удалось оформить заказ. Попробуйте позже.", &telebot.SendOptions{
			ReplyMarkup: markup,
		})
	}

	var msg string

	switch st.Code() {
	case codes.FailedPrecondition:
		msg = "❌ <b>Недостаточно средств</b>\n\nПополните баланс и попробуйте снова."
		markup.Inline(
			markup.Row(balanceBtn),
			markup.Row(backBtn),
		)

	case codes.AlreadyExists:
		msg = "ℹ️ <b>Вы уже приобрели этот товар</b>\n\nПроверьте раздел «Мои подписки»."
		markup.Inline(markup.Row(backBtn))

	case codes.Aborted:
		msg = "ℹ️ <b>Заказ был создан, но произошла внутренняя ошибка, которую мы пытаемся решить</b>\n\nПроверяйте статус заказа в разделе «Мои заказы»."
		markup.Inline(markup.Row(backBtn))

	case codes.Internal:
		msg = "⚠️ <b>Внутренняя ошибка сервера</b>\n\nПопробуйте позже."
		markup.Inline(markup.Row(backBtn))
		h.log.Error("internal error on purchase", "error", st.Message())

	default:
		msg = "⚠️ Не удалось оформить заказ. Попробуйте позже."
		markup.Inline(markup.Row(backBtn))
		h.log.Warn("unexpected grpc error", "code", st.Code(), "message", st.Message())
	}

	return ctx.Send(msg, &telebot.SendOptions{
		ParseMode:   telebot.ModeHTML,
		ReplyMarkup: markup,
	})
}

func (h *Handler) handlePurchaseSuccess(
	ctx telebot.Context,
	orderID string,
	markup *telebot.ReplyMarkup,
	backBtn telebot.Btn,
) error {
	msg := fmt.Sprintf(
		"✅ <b>Заказ успешно оформлен!</b>\n\n"+
			"Номер заказа: <code>%s</code>\n\n"+
			"Спасибо за покупку! 🎉",
		orderID,
	)

	ordersBtn := markup.Data("📦 Мои заказы", string(constants.MyOrdersCallback))
	markup.Inline(
		markup.Row(ordersBtn),
		markup.Row(backBtn),
	)

	return ctx.Send(msg, &telebot.SendOptions{
		ParseMode:   telebot.ModeHTML,
		ReplyMarkup: markup,
	})
}
