package middlewares

import (
	"github.com/GargantuaLabs/client-gateway/internal/features/auth/domain"
	"gopkg.in/telebot.v4"
)

func (m *Middlewares) OwnerOnlyMiddleware(botOwnerID domain.UserID) telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(ctx telebot.Context) error {
			user, ok := ctx.Get(CtxUserKey).(*domain.User)
			if !ok {
				return nil
			}

			if user.ID != botOwnerID {
				return nil
			}

			return next(ctx)
		}
	}
}
