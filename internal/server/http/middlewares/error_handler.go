package middlewares

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/GargantuaLabs/client-gateway/internal/features/shared"
	"github.com/gofiber/fiber/v3"
)

type ErrorHandlerMiddleware struct {
	logger *slog.Logger
}

func NewErrorHandler(logger *slog.Logger) *ErrorHandlerMiddleware {
	return &ErrorHandlerMiddleware{
		logger: logger,
	}
}

func (h *ErrorHandlerMiddleware) Handle(ctx fiber.Ctx, err error) error {
	code := http.StatusInternalServerError
	message := "Internal server error"

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
		message = e.Message
	}

	var apiErr *shared.AppError
	if errors.As(err, &apiErr) {
		code = mapKindToStatus(apiErr)
		message = apiErr.Message
	}

	if code >= http.StatusInternalServerError {
		h.logger.Error("internal server error",
			slog.Any("error", err),
			slog.String("path", ctx.Path()),
			slog.String("method", ctx.Method()),
		)

		ctx.Status(http.StatusInternalServerError)
		return ctx.JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	ctx.Status(code)
	return ctx.JSON(fiber.Map{
		"error": message,
	})
}

func mapKindToStatus(err *shared.AppError) int {
	switch err.Code {
	case shared.CodeBadRequest, shared.CodeValidation:
		return http.StatusBadRequest
	case shared.CodeAlreadyExists:
		return http.StatusConflict
	case shared.CodeUnauthorized:
		return http.StatusUnauthorized
	case shared.CodeForbidden:
		return http.StatusForbidden
	case shared.CodeNotFound:
		return http.StatusNotFound
	case shared.CodeTooManyRequests:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}
