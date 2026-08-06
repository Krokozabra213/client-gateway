package service

import (
	"context"

	authDomain "github.com/GargantuaLabs/client-gateway/internal/features/auth/domain"
	d "github.com/GargantuaLabs/client-gateway/internal/features/bot/domain"
)

func (s *Service) GetBotByUserID(ctx context.Context, userID authDomain.UserID, botID d.BotID) (*d.Bot, error) {
	return s.repo.GetBotByUserID(ctx, userID, botID)
}
