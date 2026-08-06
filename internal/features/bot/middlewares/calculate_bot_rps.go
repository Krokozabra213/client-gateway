package middlewares

import (
	"context"
	"fmt"

	authDomain "github.com/GargantuaLabs/client-gateway/internal/features/auth/domain"
	"github.com/GargantuaLabs/client-gateway/internal/features/bot/domain"
	"github.com/prometheus/client_golang/prometheus"

	"gopkg.in/telebot.v4"
)

func (m *Middlewares) CalculateBotRps(botID domain.BotID) telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(ctx telebot.Context) error {
			_, ok := ctx.Get(CtxUserKey).(*authDomain.User)
			if !ok {
				return nil
			}

			m.botRpsCounter.With(prometheus.Labels{"bot_id": botID.String()}).Inc()
			if err := m.botService.IncrementBotRps(context.Background(), botID); err != nil {
				m.log.Warn(fmt.Sprintf("Failed to increment bot rps: %v", err))
			}

			return next(ctx)
		}
	}
}
