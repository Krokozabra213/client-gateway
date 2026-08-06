package domain

import (
	"time"

	authDomain "github.com/GargantuaLabs/client-gateway/internal/features/auth/domain"
	"github.com/GargantuaLabs/client-gateway/internal/features/shared/types"
)

type Bot struct {
	ID         BotID
	CreatedAt  time.Time
	UpdatedAt  time.Time
	TelegramID types.TelegramID
	TokenHash  types.NonEmptyString
	Token      types.NonEmptyString
	Name       types.NonEmptyString
	UserID     authDomain.UserID
	IsActive   bool
}
