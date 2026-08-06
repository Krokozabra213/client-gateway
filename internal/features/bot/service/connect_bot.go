package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/GargantuaLabs/client-gateway/internal/clients"
	"github.com/GargantuaLabs/client-gateway/internal/config"
	authDomain "github.com/GargantuaLabs/client-gateway/internal/features/auth/domain"
	d "github.com/GargantuaLabs/client-gateway/internal/features/bot/domain"
	postgresRepo "github.com/GargantuaLabs/client-gateway/internal/features/bot/repository/postgres"
	"github.com/GargantuaLabs/client-gateway/internal/features/shared/types"
	"github.com/google/uuid"
	"gopkg.in/telebot.v4"
)

type ConnectBotInput struct {
	UserID   authDomain.UserID
	BotToken types.NonEmptyString
}

func (s *Service) ConnectBot(ctx context.Context, data *ConnectBotInput) (*telebot.Bot, error) {
	hashed, err := types.NewNonEmptyString(s.hashingService.Hash(data.BotToken.Value()))
	if err != nil {
		return nil, err
	}

	existingBot, err := s.repo.GetBotByTokenHash(ctx, hashed)
	if err != nil && !errors.Is(err, postgresRepo.ErrBotNotFound) {
		return nil, err
	}

	if existingBot != nil {
		return nil, ErrBotAlreadyExists
	}

	id := d.FromUuid(uuid.New())

	b, err := clients.NewBotClient(&config.BotConfig{Token: data.BotToken.Value(), WebhookUrl: fmt.Sprintf("%s/%s", s.cfg.WebhookUrl, id), Secret: s.cfg.Secret}, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating bot client: %v %w", err, ErrInvalidBot)
	}

	encrypted, err := s.encryptionService.Encrypt(data.BotToken.Value())
	if err != nil {
		return nil, err
	}

	tgID, err := types.NewTelegramID(b.Me.ID)
	if err != nil {
		return nil, ErrInvalidBot
	}

	name, err := types.NewNonEmptyString(b.Me.Username)
	if err != nil {
		return nil, ErrInvalidBot
	}

	token, err := types.NewNonEmptyString(encrypted)
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateBot(ctx, &d.Bot{
		ID:         id,
		TelegramID: tgID,
		UserID:     data.UserID,
		TokenHash:  hashed,
		Token:      token,
		Name:       name,
	}); err != nil {
		return nil, err
	}

	if err := s.cache.IncrementBotCount(ctx, data.UserID); err != nil {
		return nil, err
	}

	return b, nil
}
