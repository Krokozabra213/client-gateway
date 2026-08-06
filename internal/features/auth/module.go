package auth

import (
	"log/slog"

	"github.com/GargantuaLabs/client-gateway/internal/config"
	handler "github.com/GargantuaLabs/client-gateway/internal/features/auth/handler"
	postgresRepo "github.com/GargantuaLabs/client-gateway/internal/features/auth/repository/postgres"
	redisRepo "github.com/GargantuaLabs/client-gateway/internal/features/auth/repository/redis"
	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
	"gopkg.in/telebot.v4"

	"github.com/GargantuaLabs/client-gateway/internal/features/auth/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Module struct {
	httpHandler *handler.Handler
	Service     *service.Service
}

type Dependencies struct {
	Cfg       *config.Config
	PGXPool   *pgxpool.Pool
	Redis     *redis.Client
	BotClient *telebot.Bot
	Log       *slog.Logger
}

func New(deps *Dependencies) *Module {
	dbRepo := postgresRepo.New(deps.PGXPool)
	cacheRepo := redisRepo.New(deps.Redis)

	svc := service.New(deps.Log, cacheRepo, dbRepo, &deps.Cfg.App)
	httpHandler := handler.New(svc, deps.Log, deps.Cfg)

	return &Module{
		httpHandler: httpHandler,
		Service:     svc,
	}
}

func (m *Module) HTTPv1() func(fiber.Router) {
	return func(router fiber.Router) {
	}
}
