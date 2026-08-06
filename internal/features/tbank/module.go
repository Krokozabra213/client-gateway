package tbank

import (
	"log/slog"

	"github.com/GargantuaLabs/client-gateway/internal/clients"
	"github.com/GargantuaLabs/client-gateway/internal/features/tbank/handler"
	"github.com/gofiber/fiber/v3"
)

type Dependencies struct {
	Payment *clients.PaymentClient
	Log     *slog.Logger
}

type Module struct {
	handler *handler.Handler
	logger  *slog.Logger
}

func (m *Module) HTTPv1() func(fiber.Router) {
	return func(r fiber.Router) {
		r.Post("/notification", m.handler.HandleNotification)
	}
}

func New(deps *Dependencies) *Module {

    h := handler.New(deps.Log, deps.Payment)

	return &Module{
		logger: deps.Log,
        handler: h,
	}
}
