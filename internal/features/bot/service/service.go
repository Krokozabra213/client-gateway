package service

import (
	"context"

	"github.com/GargantuaLabs/client-gateway/internal/config"
	authDomain "github.com/GargantuaLabs/client-gateway/internal/features/auth/domain"
	"github.com/GargantuaLabs/client-gateway/internal/features/auth/service"
	d "github.com/GargantuaLabs/client-gateway/internal/features/bot/domain"
	"github.com/GargantuaLabs/client-gateway/internal/features/shared/types"
)

type AuthService interface {
	Register(ctx context.Context, data *service.RegisterInput) error
	GetUserByTelegramID(ctx context.Context, telegramID types.TelegramID) (*authDomain.User, error)
}

type EncryptionService interface {
	Encrypt(plainText string) (string, error)
	Decrypt(plainText string) (string, error)
}

type HashingService interface {
	Hash(plainText string) string
}

type Cache interface {
	StoreBotCount(ctx context.Context, count int, ownerID authDomain.UserID) error
	GetBotCount(ctx context.Context, ownerID authDomain.UserID) (int, error)
	IncrementBotCount(ctx context.Context, ownerID authDomain.UserID) error
	StoreBotData(ctx context.Context, botID d.BotID, botToken types.NonEmptyString, ownerID authDomain.UserID) error
	GetBotData(ctx context.Context, botID d.BotID) (botToken types.NonEmptyString, ownerID authDomain.UserID, err error)
	IncrementBotRps(ctx context.Context, botID d.BotID) error
	GetBotLastMinuteRps(ctx context.Context, botID d.BotID) (min, avg, max float64, err error)
}

type Repository interface {
	CreateBot(ctx context.Context, bot *d.Bot) error
	GetBotsByUserID(ctx context.Context, userID authDomain.UserID, pagination *types.Pagination) ([]*d.Bot, error)
	GetBotByTokenHash(ctx context.Context, token types.NonEmptyString) (*d.Bot, error)
	GetBotByID(ctx context.Context, id d.BotID) (*d.Bot, error)
	GetBotByUserID(ctx context.Context, userID authDomain.UserID, botID d.BotID) (*d.Bot, error)
	GetBotCountByUserID(ctx context.Context, userID authDomain.UserID) (int, error)
}

// type botEntry struct {
// 	bot        *telebot.Bot
// 	ownerID    authDomain.UserID
// 	expiration time.Time
// }

type Service struct {
	authService       AuthService
	encryptionService EncryptionService
	hashingService    HashingService

	cfg   *config.BotConfig
	repo  Repository
	cache Cache
}

func New(authService AuthService, encryptionService EncryptionService, hashingService HashingService, repo Repository, cache Cache, cfg *config.BotConfig) *Service {
	svc := &Service{
		authService:       authService,
		encryptionService: encryptionService,
		hashingService:    hashingService,
		cfg:               cfg,
		repo:              repo,
		cache:             cache,
	}

	return svc
}
