package handler

import (
	"context"
	"strings"
	"time"

	authDomain "github.com/GargantuaLabs/client-gateway/internal/features/auth/domain"
	"github.com/GargantuaLabs/client-gateway/internal/features/bot/constants"
	"github.com/GargantuaLabs/client-gateway/internal/features/shared"
	"github.com/GargantuaLabs/client-gateway/internal/features/shared/types"
	"gopkg.in/telebot.v4"
)

const birthDateSeparator = ": "

func (h *Handler) registerCallback(ctx telebot.Context) error {
	if ctx.Callback() != nil {
		defer ctx.Respond() //nolint:errcheck
	}

	if err := ctx.Edit(h.buildStartupMessage()); err != nil {
		h.log.Warn(err.Error())
	}

	return ctx.Send(h.buildRegisterMessage(), h.markups[constants.RegisterMarkup])
}

func (h *Handler) registerCancelCallback(ctx telebot.Context) error {
	if ctx.Callback() != nil {
		defer ctx.Respond() //nolint:errcheck
	}

	telegramID, err := types.NewTelegramID(ctx.Sender().ID)
	if err != nil {
		return ctx.Send("Некорректный телеграмм ID.")
	}

	if err := ctx.Edit(h.buildRegisterMessage()); err != nil {
		h.log.Warn(err.Error())
	}

	if err := h.sendAndStoreMessage(ctx, telegramID, "🚫 Регистрация прекращена. Для продолжения регистрации введи команду: /start", nil); err != nil {
		return shared.NewAppError(shared.CodeInternal, "failed sendAndStoreMessage", err)
	}

	return nil
}

func (h *Handler) registerApproveCallback(ctx telebot.Context) error {
	if ctx.Callback() != nil {
		defer ctx.Respond() //nolint:errcheck
	}

	telegramID, err := types.NewTelegramID(ctx.Sender().ID)
	if err != nil {
		return ctx.Send("Некорректный телеграмм ID.")
	}

	if err := ctx.Edit(h.buildRegisterMessage()); err != nil {
		h.log.Warn(err.Error())
	}

	if err := h.sendAndStoreMessage(ctx, telegramID, `
👤 Укажи свой пол:

<i>*После регистрации пол нельзя будет сменить самостоятельно. Выбирай с умом</i>`, nil, h.markups[constants.GenderMarkup]); err != nil {
		return shared.NewAppError(shared.CodeInternal, "failed sendAndStoreMessage", err)
	}

	return nil
}

func (h *Handler) genderCallback(ctx telebot.Context) error {
	if ctx.Callback() != nil {
		defer ctx.Respond() //nolint:errcheck
	}

	telegramID, err := types.NewTelegramID(ctx.Sender().ID)
	if err != nil {
		return ctx.Send("Некорректный телеграмм ID.")
	}

	parsedGender, err := authDomain.NewUserGender(ctx.Callback().Unique)
	if err != nil {
		return ctx.Send("Некорректный пол.")
	}

	if err := h.cache.StoreGender(context.Background(), telegramID, parsedGender); err != nil {
		return shared.NewAppError(shared.CodeInternal, "failed StoreGender", err)
	}

	if err := h.removeStoredMessage(ctx, telegramID); err != nil {
		return shared.NewAppError(shared.CodeInternal, "failed removeStoredMessage", err)
	}

	if err := h.sendAndStoreMessage(ctx, telegramID, h.buildBirthdayMessage("XX.XX.XXXX"), nil, h.markups[constants.BirthdayMarkup]); err != nil {
		return shared.NewAppError(shared.CodeInternal, "failed sendAndStoreMessage", err)
	}

	return nil
}

func (h *Handler) birthdayCallback(ctx telebot.Context) error {
	if ctx.Callback() != nil {
		defer ctx.Respond() //nolint:errcheck
	}

	text := ctx.Text()
	parts := strings.Split(text, birthDateSeparator)

	if len(parts) < 2 {
		h.log.Warn("birthdayCallback: invalid mask state", "text", text)
		return ctx.Send("⚠️ Возникла проблема. Попробуйте позже.")
	}

	maskLength := 10
	datePart := parts[1]

	if len(datePart) < maskLength {
		h.log.Warn("birthdayCallback: invalid mask length", "text", text)
		return ctx.Send("⚠️ Возникла проблема. Попробуйте позже.")
	}

	datePart = parts[1][:maskLength]
	indexOfNotFilled := strings.Index(datePart, "X")
	action := ctx.Callback().Unique

	if action == "b_ok" {
		if indexOfNotFilled != -1 {
			return nil
		}

		if err := h.register(ctx, datePart); err != nil {
			return err
		}

		return nil
	}

	out := []rune(datePart)

	if action == "b_del" {
		i := strings.LastIndexFunc(datePart, func(r rune) bool { return r != 'X' && r != '.' })
		if i == -1 {
			return nil
		}

		out[i] = 'X'

		if err := ctx.Edit(
			h.buildBirthdayMessage(string(out)),
			h.markups[constants.BirthdayMarkup],
		); err != nil {
			return err
		}

		return nil
	}

	if indexOfNotFilled == -1 {
		return nil
	}

	r := rune(ctx.Data()[0])

	if !h.isValidBirthDate(&out, indexOfNotFilled, int(r-'0')) {
		return ctx.Respond(&telebot.CallbackResponse{
			ShowAlert: true,
			Text:      "⚠️ Некорректная дата!",
		})
	}

	out[indexOfNotFilled] = r

	if err := ctx.Edit(
		h.buildBirthdayMessage(string(out)),
		h.markups[constants.BirthdayMarkup],
	); err != nil {
		return err
	}

	return nil
}

func (h *Handler) isValidBirthDate(out *[]rune, index int, digit int) bool {
	if index == 0 && digit > 3 {
		return false
	}

	if index == 1 && int((*out)[0]-'0') == 3 && digit > 1 {
		return false
	}

	if index == 3 && digit > 1 {
		return false
	}

	if index == 4 && int((*out)[3]-'0') == 1 && digit > 2 {
		return false
	}

	currYear := time.Now().Year()
	currYDigits := []int{
		currYear / 1000,
		(currYear / 100) % 10,
		(currYear / 10) % 10,
		currYear % 10,
	}

	if index == 6 && (digit < 1 || digit > currYDigits[0]) {
		return false
	}

	if index == 7 {
		firstDigit := int((*out)[6] - '0')
		if firstDigit == currYDigits[0] && digit > currYDigits[1] {
			return false
		}
	}

	if index == 8 {
		yearSoFar := int((*out)[6]-'0')*10 + int((*out)[7]-'0')
		if yearSoFar == currYDigits[0]*10+currYDigits[1] && digit > currYDigits[2] {
			return false
		}
	}

	if index == 9 {
		yearSoFar := int((*out)[6]-'0')*100 + int((*out)[7]-'0')*10 + int((*out)[8]-'0')
		if yearSoFar == currYDigits[0]*100+currYDigits[1]*10+currYDigits[2] && digit > currYDigits[3] {
			return false
		}
	}

	return true
}
