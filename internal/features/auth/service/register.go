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

type RegisterInput struct {
	TelegramID        types.TelegramID
	TelegramFirstName d.Username
	BirthDate         time.Time
	Gender            d.UserGender
}

func (s *Service) Register(ctx context.Context, data *RegisterInput) error {
	user, err := s.repo.GetUserByTelegramID(ctx, data.TelegramID)
	if err != nil && !errors.Is(err, postgresRepo.ErrUserNotFound) {
		return shared.NewAppError(shared.CodeInternal, "db is falling...", err)
	}

	if user != nil {
		return shared.NewAppError(shared.CodeAlreadyExists, "Вы уже зарегистрированы", nil)
	}

	_, err = s.repo.CreateUser(ctx, d.User{
		TelegramID:        data.TelegramID,
		BirthDate:         data.BirthDate,
		Gender:            data.Gender,
		TelegramFirstName: data.TelegramFirstName,
		IsVerified:        true,
		Role:              d.RoleRegular,
	})

	return err
}
