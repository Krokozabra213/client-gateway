package handler

import (
	"context"
	"fmt"

	authDomain "github.com/GargantuaLabs/client-gateway/internal/features/auth/domain"
	"github.com/GargantuaLabs/client-gateway/internal/features/bot/constants"
	"github.com/GargantuaLabs/client-gateway/internal/features/bot/middlewares"
	"github.com/GargantuaLabs/client-gateway/internal/features/shared"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/telebot.v4"
)

func (h *Handler) balanceCallback(ctx telebot.Context) error {
	if ctx.Callback() != nil {
		_ = ctx.Respond() //nolint:errcheck
	}

	user, ok := ctx.Get(middlewares.CtxUserKey).(*authDomain.User)
	if !ok {
		return ctx.Send("Не удалось определить пользователя. Попробуйте заново.")
	}

	var currentBalance int64
	resp, err := h.balance.GetAccount(context.Background(), user.ID.String())
	if err != nil {
		st, ok := status.FromError(err)
		if !ok || st.Code() != codes.NotFound {
			return shared.NewAppError(shared.CodeBadRequest, "Ошибка получения баланса. Попробуйте чуть позже.", err)
		}
	} else {
		currentBalance = resp.GetBalance()
	}

	msg := fmt.Sprintf(`<b>Информация о балансе</b>

💳 Текущий баланс %.2f ₽
`, float64(currentBalance)/100)

	if currentBalance <= 0 {
		msg += "\n❗️ Баланс пуст — сервисы остановлены"
	}

	return ctx.Send(msg, h.markups[constants.BalanceMarkup])
}
