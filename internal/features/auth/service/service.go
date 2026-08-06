package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/GargantuaLabs/client-gateway/internal/config"
	d "github.com/GargantuaLabs/client-gateway/internal/features/auth/domain"
	"github.com/GargantuaLabs/client-gateway/internal/features/shared/types"
)

type Cache interface {
	GetUser(ctx context.Context, telegramID types.TelegramID) (*d.User, error)
	SetUser(ctx context.Context, user *d.User, ttl time.Duration) error
}

type Repository interface {
	CreateUser(ctx context.Context, user d.User) (d.UserID, error)
	GetUserByID(ctx context.Context, id d.UserID) (*d.User, error)
	GetUserByTelegramID(ctx context.Context, telegramID types.TelegramID) (*d.User, error)
}

type Service struct {
	log    *slog.Logger
	cache  Cache
	repo   Repository
	config *config.AppConfig
}

func New(log *slog.Logger, cache Cache, repo Repository, config *config.AppConfig) *Service {
	return &Service{
		log,
		cache,
		repo,
		config,
	}
}
