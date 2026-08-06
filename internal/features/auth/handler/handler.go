package handler

import (
	"log/slog"

	"github.com/GargantuaLabs/client-gateway/internal/config"
	"github.com/GargantuaLabs/client-gateway/internal/features/auth/service"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	cfg     *config.Config
	log     *slog.Logger
	service *service.Service
}

func New(svc *service.Service, log *slog.Logger, cfg *config.Config) *Handler {
	return &Handler{
		cfg:     cfg,
		log:     log,
		service: svc,
	}
}

func (h *Handler) TestAuth(c fiber.Ctx) error {
	return c.JSON(fiber.Map{"test": "ok"})
}
