package middlewares

import (
	"context"
	"log/slog"

	d "github.com/GargantuaLabs/client-gateway/internal/features/auth/domain"
	"github.com/GargantuaLabs/client-gateway/internal/features/bot/domain"
	"github.com/GargantuaLabs/client-gateway/internal/features/shared/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type AuthService interface {
	GetUserByTelegramID(ctx context.Context, telegramID types.TelegramID) (*d.User, error)
}

type BotService interface {
	IncrementBotRps(ctx context.Context, botID domain.BotID) error
}

type Middlewares struct {
	authService   AuthService
	botService    BotService
	log           *slog.Logger
	botRpsCounter *prometheus.CounterVec
}

func New(authService AuthService, botService BotService, registerer prometheus.Registerer, log *slog.Logger) *Middlewares {
	return &Middlewares{
		authService: authService,
		botService:  botService,
		log:         log,
		botRpsCounter: promauto.With(registerer).NewCounterVec(
			prometheus.CounterOpts{
				Name: "bot_rps_total",
				Help: "Calculate client bot RPS",
			}, []string{"bot_id"}),
	}
}
