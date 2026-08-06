package service

import (
	"context"

	authDomain "github.com/GargantuaLabs/client-gateway/internal/features/auth/domain"
	d "github.com/GargantuaLabs/client-gateway/internal/features/bot/domain"
	"github.com/GargantuaLabs/client-gateway/internal/features/shared/types"
)

func (s *Service) GetBotsByUserID(ctx context.Context, userID authDomain.UserID, pagination *types.Pagination) ([]*d.Bot, int, error) {
	bots, err := s.repo.GetBotsByUserID(ctx, userID, pagination)
	if err != nil {
		return nil, 0, err
	}

	botCount, err := s.GetBotCount(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	return bots, botCount, nil
}
