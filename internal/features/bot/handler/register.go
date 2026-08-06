package handler

import (
	"context"
	"fmt"
	"time"

	authDomain "github.com/GargantuaLabs/client-gateway/internal/features/auth/domain"
	"github.com/GargantuaLabs/client-gateway/internal/features/auth/service"
	"github.com/GargantuaLabs/client-gateway/internal/features/bot/constants"
	"github.com/GargantuaLabs/client-gateway/internal/features/shared/types"
	"gopkg.in/telebot.v4"
)

func (h *Handler) register(ctx telebot.Context, dateOfBirth string) error {
	parsedDate, err := time.Parse("02.01.2006", dateOfBirth)
	if err != nil {
		return ctx.Send("Некорректная дата.")
	}

	telegramID, err := types.NewTelegramID(ctx.Sender().ID)
	if err != nil {
		return ctx.Send("Некорректный телеграмм ID.")
	}

	firstName, err := authDomain.NewUsername(ctx.Sender().FirstName)
	if err != nil {
		return ctx.Send("Некорректный first name")
	}

	gender, err := h.cache.GetAndDeleteGender(context.Background(), telegramID)
	if err != nil {
		return err
	}

	if err := h.authService.Register(context.Background(), &service.RegisterInput{
		TelegramID:        telegramID,
		BirthDate:         parsedDate,
		Gender:            gender,
		TelegramFirstName: firstName,
	}); err != nil {
		return err
	}

	if err := h.removeStoredMessage(ctx, telegramID); err != nil {
		return err
	}

	var label string
	if gender == authDomain.GenderMale {
		label = "прошел"
	} else {
		label = "прошла"
	}

	return ctx.Send(fmt.Sprintf(`🎉 Поздравляю! Ты %s регистрацию. Теперь ты можешь перейти в главное меню или ознакомиться с функционалом`, label), h.markups[constants.RegisterSuccessMarkup])
}
