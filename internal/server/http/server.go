package httpserver

import (
	"context"
	"log/slog"

	"github.com/GargantuaLabs/client-gateway/internal/server/http/middlewares"
	"github.com/Krokozabra213/gargantua_common/pkg/config"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/google/uuid"
)

type Server struct {
	App    *fiber.App
	cfg    *config.HTTPConfig
	logger *slog.Logger
}

func New(cfg *config.HTTPConfig, logger *slog.Logger) *Server {
	errorHandler := middlewares.NewErrorHandler(logger)

	fiberApp := fiber.New(fiber.Config{
		ReadTimeout:    cfg.ReadTimeout,
		WriteTimeout:   cfg.WriteTimeout,
		ReadBufferSize: cfg.MaxHeaderBytes * 1024 * 1024,
		BodyLimit:      cfg.BodyLimitBytes * 1024 * 1024,
		ErrorHandler:   errorHandler.Handle,
	})

	fiberApp.Use(
		cors.New(cors.Config{
			AllowOrigins:     cfg.AllowedOrigins,
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
			AllowCredentials: true,
		}),
		recover.New(),
		requestid.New(requestid.Config{
			Header: "X-Request-Id",
			Generator: func() string {
				return uuid.New().String()
			},
		}),
		middlewares.RequestLogger(logger),
	)

	return &Server{
		App:    fiberApp,
		cfg:    cfg,
		logger: logger,
	}
}

func (a *Server) Run() error {
	a.logger.Info("server starting", "addr", a.cfg.HTTPAddress())
	return a.App.Listen(a.cfg.HTTPAddress())
}

func (a *Server) Shutdown(ctx context.Context) error {
	a.logger.Info("shutting down server...")
	return a.App.ShutdownWithContext(ctx)
}
