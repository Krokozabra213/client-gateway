package errorhandler

import (
	"errors"
	"log/slog"

	"github.com/GargantuaLabs/client-gateway/internal/features/shared"
	"gopkg.in/telebot.v4"
)

type BotErrorHandler struct {
	logger *slog.Logger
}

func NewBotErrorHandler(logger *slog.Logger) *BotErrorHandler {
	return &BotErrorHandler{
		logger: logger,
	}
}

func (h *BotErrorHandler) Handle(err error, ctx telebot.Context) {
	message := "Что-то пошло не так, попробуйте чуть позже"

	var userID int64
	if ctx != nil && ctx.Sender() != nil {
		userID = ctx.Sender().ID
	}

	var apiErr *shared.AppError
	if errors.As(err, &apiErr) {
		if apiErr.Code != shared.CodeInternal {
			message = apiErr.Message

			h.logger.Warn("bot app error",
				slog.String("code", string(apiErr.Code)),
				slog.String("message", apiErr.Message),
				slog.String("error", apiErr.Error()),
				slog.Int64("user_id", userID),
			)
		} else {
			h.logger.Error("bot internal app error",
				slog.String("code", string(apiErr.Code)),
				slog.String("error", apiErr.Error()),
				slog.Int64("user_id", userID),
			)
		}
	} else {
		h.logger.Error("bot unexpected error",
			slog.Any("error", err),
			slog.Int64("user_id", userID),
		)
	}

	if ctx != nil {
		_ = ctx.Send(message)
	}
}
