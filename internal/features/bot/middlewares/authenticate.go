package middlewares

import (
	"context"

	"github.com/GargantuaLabs/client-gateway/internal/features/shared"
	"github.com/GargantuaLabs/client-gateway/internal/features/shared/types"

	"gopkg.in/telebot.v4"
)

const CtxUserKey = "user"

func (m *Middlewares) Authenticate(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		telegramID, err := types.NewTelegramID(ctx.Sender().ID)
		if err != nil {
			return ctx.Send("Некорректный телеграмм ID.")
		}

		user, err := m.authService.GetUserByTelegramID(context.Background(), telegramID)
		if err != nil {
			return err
		}

		if user == nil {
			return shared.NewAppError(shared.CodeNotFound, "Вы не авторизованы.", err)
		}

		ctx.Set(CtxUserKey, user)

		return next(ctx)
	}
}
