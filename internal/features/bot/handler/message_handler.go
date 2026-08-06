package handler

import (
	"context"

	d "github.com/GargantuaLabs/client-gateway/internal/features/auth/domain"
	"github.com/GargantuaLabs/client-gateway/internal/features/bot/constants"
	"github.com/GargantuaLabs/client-gateway/internal/features/bot/middlewares"
	"gopkg.in/telebot.v4"
)

func (h *Handler) HandleMessage(ctx telebot.Context) error {
	user, ok := ctx.Get(middlewares.CtxUserKey).(*d.User)
	if !ok {
		return nil
	}

	lastAction, err := h.cache.GetLastAction(context.Background(), user.TelegramID)
	if err != nil {
		h.log.Warn(err.Error())
	}

	switch lastAction {
	case constants.ConnectBotAction:
		return h.connectBot(ctx)
	//case constants.DepositTbankAction:
	//	return h.depositTbank(ctx)

	default:
		return nil
	}
}
