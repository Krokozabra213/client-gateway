package service

import (
	"context"

	authDomain "github.com/GargantuaLabs/client-gateway/internal/features/auth/domain"
	d "github.com/GargantuaLabs/client-gateway/internal/features/bot/domain"
	"github.com/google/uuid"
	"gopkg.in/telebot.v4"
)

func (s *Service) GetConnectedBot(ctx context.Context, botID d.BotID) (*telebot.Bot, authDomain.UserID, error) {
	token, userID, err := s.cache.GetBotData(ctx, botID)
	if err != nil {
		return nil, authDomain.UserID(uuid.Nil), err
	}

	if token.Value() == "" {
		botData, err := s.repo.GetBotByID(ctx, botID)
		if err != nil {
			return nil, authDomain.UserID(uuid.Nil), err
		}

		if err := s.cache.StoreBotData(ctx, botID, botData.Token, botData.UserID); err != nil {
			return nil, authDomain.UserID(uuid.Nil), err
		}

		token = botData.Token
		userID = botData.UserID
	}

	decrypted, err := s.encryptionService.Decrypt(token.Value())
	if err != nil {
		return nil, authDomain.UserID(uuid.Nil), err
	}

	bot, err := telebot.NewBot(telebot.Settings{
		Token:     decrypted,
		ParseMode: telebot.ModeHTML,
	})
	if err != nil {
		return nil, authDomain.UserID(uuid.Nil), err
	}

	return bot, userID, nil
}
