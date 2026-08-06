package service

import (
	"context"

	authDomain "github.com/GargantuaLabs/client-gateway/internal/features/auth/domain"
)

func (s *Service) GetBotCount(ctx context.Context, userID authDomain.UserID) (int, error) {
	cached, err := s.cache.GetBotCount(ctx, userID)
	if err != nil {
		return 0, err
	}

	if cached == 0 {
		count, err := s.repo.GetBotCountByUserID(ctx, userID)
		if err != nil {
			return 0, err
		}

		if err := s.cache.StoreBotCount(ctx, count, userID); err != nil {
			return 0, err
		}

		return count, nil
	}

	return cached, nil
}
