package domain

import (
	"time"

	"github.com/GargantuaLabs/client-gateway/internal/features/shared/types"
)

type User struct {
	ID                UserID
	CreatedAt         time.Time
	UpdatedAt         time.Time
	TelegramID        types.TelegramID
	TelegramFirstName Username
	Username          *Username
	BirthDate         time.Time
	Gender            UserGender
	Role              UserRole
	IsBanned          bool
	IsVerified        bool
}
