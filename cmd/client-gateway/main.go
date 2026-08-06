package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GargantuaLabs/client-gateway/internal/clients"
	"github.com/GargantuaLabs/client-gateway/internal/config"
	"github.com/GargantuaLabs/client-gateway/internal/features/auth"
	"github.com/GargantuaLabs/client-gateway/internal/features/bot"
	"github.com/GargantuaLabs/client-gateway/internal/features/tbank"
	"github.com/GargantuaLabs/client-gateway/internal/security"
	httpserver "github.com/GargantuaLabs/client-gateway/internal/server/http"
	commonClients "github.com/Krokozabra213/gargantua_common/pkg/clients"
	"github.com/Krokozabra213/gargantua_common/pkg/logger"
	"github.com/gofiber/fiber/v3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const shutdown = 5 * time.Second

func main() {
	if err := run(); err != nil {
		slog.Error("application failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Init()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// otelHandler, cleanup, err := telemetry.NewBuilder().
	// 	WithServiceName("client-gateway").
	// 	WithOTELEndpoint("otel-collector:4317").
	// 	WithLogs().
	// 	WithMetrics().
	// 	WithTraces().
	// 	Build(ctx)

	// log := logger.Init(&cfg.Logger, otelHandler)
	log := logger.Init(&cfg.Logger)

	name := "go.opentelemetry.io/contrib/examples/otel-collector"
	tracer := otel.Tracer(name)
	meter := otel.Meter(name)

	commonAttrs := []attribute.KeyValue{
		attribute.String("attrA", "chocolate"),
		attribute.String("attrB", "raspberry"),
		attribute.String("attrC", "vanilla"),
	}

	runCount, err := meter.Int64Counter("run", metric.WithDescription("The number of times the iteration ran"))
	if err != nil {
		log.Error("message", "error", err)
	}

	ctx, span := tracer.Start(
		ctx,
		"CollectorExporter-Example",
		trace.WithAttributes(commonAttrs...))
	defer span.End()
	for i := range 10 {
		_, iSpan := tracer.Start(ctx, fmt.Sprintf("Sample-%d", i))
		runCount.Add(ctx, 1, metric.WithAttributes(commonAttrs...))
		log.Info("Doing really hard work (%d / 10)\n", "count", "i+1")

		<-time.After(time.Second)
		iSpan.End()
	}

	reg := prometheus.DefaultRegisterer

	botClient, err := clients.NewBotClient(&cfg.Bot, log)
	if err != nil {
		panic(err)
	}

	pool, err := commonClients.NewPostgresClient(&cfg.Postgres)
	if err != nil {
		panic(err)
	}

	paymentClient, err := clients.NewPayment(ctx, cfg.Payment.Addr)
	if err != nil {
		return err
	}
	defer func() {
		if err := paymentClient.Close(); err != nil {
			log.Error("ошибка при закрытии paymentClient: %v", "error", err)
		}
	}()

	balanceClient, err := clients.NewBalance(ctx, cfg.Balance.Addr)
	if err != nil {
		return err
	}
	defer func() {
		if err := balanceClient.Close(); err != nil {
			log.Error("ошибка при закрытии balanceClient: %v", "error", err)
		}
	}()

	orderProductClient, err := clients.NewOrderProduct(ctx, cfg.OrderProduct.Addr)
	if err != nil {
		return err
	}
	defer func() {
		if err := orderProductClient.Close(); err != nil {
			log.Error("ошибка при закрытии orderProductClient: %v", "error", err)
		}
	}()

	redisClient, err := clients.NewRedisClient(&cfg.Redis)
	if err != nil {
		panic(err)
	}

	encryptionService, _ := security.NewEncryptionService(cfg.App.EncryptionKey)
	hashingService := security.NewHashingService(cfg.App.EncryptionKey)

	authModule := auth.New(&auth.Dependencies{
		Cfg: cfg, PGXPool: pool, Redis: redisClient, BotClient: botClient, Log: log,
	})
	botModule := bot.New(&bot.Dependencies{
		BotCfg:            &cfg.Bot,
		BotClient:         botClient,
		OrderProduct:      orderProductClient,
		Redis:             redisClient,
		AuthService:       authModule.Service,
		EncryptionService: encryptionService,
		HashingService:    hashingService,
		PGXPool:           pool,
		Log:               log,
		Registerer:        reg,
		Payment:           paymentClient,
		Balance:           balanceClient,
	})
	tbankModule := tbank.New(&tbank.Dependencies{
		Payment: paymentClient, Log: log,
	})

	server := httpserver.New(&cfg.HTTP, log)
	server.App.Get("/metrics", promhttp.Handler())
	server.App.Get("/health", func(c fiber.Ctx) error {
		ctx := c.Context()
		status := fiber.Map{"status": "ok"}

		if err := redisClient.Ping(ctx).Err(); err != nil {
			status["status"] = "error"
			status["redis"] = "down"
			c.Status(http.StatusInternalServerError)
		}

		if err := pool.Ping(ctx); err != nil {
			status["status"] = "error"
			status["postgres"] = "down"
			c.Status(http.StatusInternalServerError)
		}

		return c.JSON(status)
	})

	v1 := server.App.Group("/v1")
	v1.Route("/auth", authModule.HTTPv1())
	v1.Route("/bot", botModule.HTTPv1())
	v1.Route("/tbank", tbankModule.HTTPv1())

	errCh := make(chan error, 1)
	go func() {
		log.Info("server started", "address", cfg.HTTP.HTTPAddress())
		if err := server.Run(); !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case sig := <-ctx.Done():
		log.Info("received shutdown signal", "signal", sig)
	case err := <-errCh:
		log.Error("server error", "error", err)
		return err
	}

	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdown)
	defer cancel()

	// if err := cleanup(ctx); err != nil {
	// 	log.Error("telemetry cleanup error", "error", err)
	// }

	if err := server.Shutdown(ctx); err != nil {
		log.Error("server shutdown error", "error", err)
	}

	if err := redisClient.Close(); err != nil {
		log.Error("redis client close error", "error", err)
	}

	pool.Close()

	log.Info("application soft stopped")

	return nil
}
