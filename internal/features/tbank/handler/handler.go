package handler

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/GargantuaLabs/client-gateway/internal/clients"
	"github.com/GargantuaLabs/client-gateway/internal/domain"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	payment *clients.PaymentClient
	log     *slog.Logger
}

func New(log *slog.Logger, payment *clients.PaymentClient) *Handler {
	handler := &Handler{
		payment: payment,
		log:     log,
	}

	return handler
}

func (h *Handler) HandleNotification(c fiber.Ctx) error {
	var req domain.TBankPaymentNotification

	if err := json.Unmarshal(c.Body(), &req); err != nil {
		h.log.Error("failed to parse t-bank webhook", slog.Any("err", err))
		return c.Status(fiber.StatusBadRequest).SendString("Bad Request")
	}

	err := h.payment.TbankWebhook(context.Background(), &req)
	if err != nil {
		h.log.Error("failed to handle tbank webhook", slog.Any("err", err))
		return c.Status(fiber.StatusBadRequest).SendString("Bad Request")
	}

	return c.SendString("OK")
}
