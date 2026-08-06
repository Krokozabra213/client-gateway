package handler

import (
	"context"
	"errors"
	"fmt"

	d "github.com/GargantuaLabs/client-gateway/internal/features/auth/domain"
	postgresRepo "github.com/GargantuaLabs/client-gateway/internal/features/auth/repository/postgres"
	"github.com/GargantuaLabs/client-gateway/internal/features/bot/constants"
	"github.com/GargantuaLabs/client-gateway/internal/features/bot/middlewares"
	"github.com/GargantuaLabs/client-gateway/internal/features/shared/types"
	"gopkg.in/telebot.v4"
)

func (h *Handler) Start(ctx telebot.Context) error {
	telegramID, err := types.NewTelegramID(ctx.Sender().ID)
	if err != nil {
		return ctx.Send("Некорректный телеграмм ID.")
	}

	if err := h.removeStoredMessage(ctx, telegramID); err != nil {
		h.log.Warn(err.Error())
	}

	if err := h.cache.RemoveLastAction(context.Background(), telegramID); err != nil {
		h.log.Warn(err.Error())
	}

	user, err := h.authService.GetUserByTelegramID(context.Background(), telegramID)
	if err != nil && !errors.Is(err, postgresRepo.ErrUserNotFound) {
		return err
	}

	if user != nil {
		var label string
		if user.Gender == d.GenderMale {
			label = "прошел"
		} else {
			label = "прошла"
		}

		return ctx.Send(fmt.Sprintf("🐱 Ты уже %s регистрацию! Чтобы воспользоваться функционалом платформы или ознакомиться с ним, нажми соответствующую кнопку ниже 👇", label), h.markups[constants.RegisterSuccessMarkup])
	}

	return ctx.Send(h.buildStartupMessage(), h.markups[constants.StartupMarkup])
}

func (h *Handler) Menu(ctx telebot.Context) error {
	user, ok := ctx.Get(middlewares.CtxUserKey).(*d.User)
	if !ok {
		return nil
	}

	message, markup, err := h.buildMainMenu(user)
	if err != nil {
		return err
	}

	return ctx.Send(message, markup)
}
