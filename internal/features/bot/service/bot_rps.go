package service

import (
	"context"

	"github.com/GargantuaLabs/client-gateway/internal/features/bot/domain"
)

func (s *Service) IncrementBotRps(ctx context.Context, botID domain.BotID) error {
	return s.cache.IncrementBotRps(ctx, botID)
}

func (s *Service) GetBotLastMinuteRps(ctx context.Context, botID domain.BotID) (min, avg, max float64, err error) {
	return s.cache.GetBotLastMinuteRps(ctx, botID)
}
