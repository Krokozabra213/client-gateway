package service

import (
	"context"
	"errors"
	"time"

	d "github.com/GargantuaLabs/client-gateway/internal/features/auth/domain"
	postgresRepo "github.com/GargantuaLabs/client-gateway/internal/features/auth/repository/postgres"
	"github.com/GargantuaLabs/client-gateway/internal/features/shared"
	"github.com/GargantuaLabs/client-gateway/internal/features/shared/types"
)

func (s *Service) GetUserByTelegramID(ctx context.Context, telegramID types.TelegramID) (*d.User, error) {
	userFromCache, err := s.cache.GetUser(ctx, telegramID)
	if err != nil {
		return nil, shared.NewAppError(shared.CodeInternal, "failed get user from cache", err)
	}

	if userFromCache != nil {
		return userFromCache, nil
	}

	userFromDB, err := s.repo.GetUserByTelegramID(ctx, telegramID)
	if err != nil {
		if errors.Is(err, postgresRepo.ErrUserNotFound) {
			return nil, shared.NewAppError(shared.CodeNotFound, "Вы не прошли авторизацию.", err)
		}
		return nil, shared.NewAppError(shared.CodeInternal, "failed get user from db", err)
	}

	if userFromDB != nil {
		if err := s.cache.SetUser(ctx, userFromDB, time.Minute*time.Duration(s.config.UserCacheTTLMinutes)); err != nil {
			s.log.Error("failed setUser in cache", "err", err.Error())
		}
	}

	return userFromDB, nil
}
