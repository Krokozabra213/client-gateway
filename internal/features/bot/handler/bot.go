package handler

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	authDomain "github.com/GargantuaLabs/client-gateway/internal/features/auth/domain"
	d "github.com/GargantuaLabs/client-gateway/internal/features/bot/domain"
	"github.com/GargantuaLabs/client-gateway/internal/features/bot/service"
	"github.com/GargantuaLabs/client-gateway/internal/features/shared"
	"github.com/GargantuaLabs/client-gateway/internal/features/shared/types"
	"gopkg.in/telebot.v4"

	"github.com/GargantuaLabs/client-gateway/internal/features/bot/constants"
	"github.com/GargantuaLabs/client-gateway/internal/features/bot/middlewares"
)

const (
	botsLimit          int = 2
	botNumberSeparator     = "_"
)

func (h *Handler) connectBotCallback(ctx telebot.Context) error {
	if ctx.Callback() != nil {
		defer ctx.Respond() //nolint:errcheck
	}

	if _, ok := ctx.Get(middlewares.CtxUserKey).(*authDomain.User); !ok {
		return nil
	}

	return ctx.Send(h.buildConnectBotMessage(), h.markups[constants.ConnectBotConfirmMarkup])
}

func (h *Handler) connectBotConfirmCallback(ctx telebot.Context) error {
	if ctx.Callback() != nil {
		defer ctx.Respond() //nolint:errcheck
	}

	user, ok := ctx.Get(middlewares.CtxUserKey).(*authDomain.User)
	if !ok {
		return nil
	}

	if err := h.cache.StoreLastAction(context.Background(), constants.ConnectBotAction, user.TelegramID); err != nil {
		return err
	}

	return ctx.Send("🔑 Отправьте в этот чат токен своего бота и ничего больше:")
}

func (h *Handler) connectBot(ctx telebot.Context) error {
	user, ok := ctx.Get(middlewares.CtxUserKey).(*authDomain.User)
	if !ok {
		return nil
	}

	token, err := types.NewNonEmptyString(ctx.Text())
	if err != nil {
		return ctx.Send("Токен бота не может быть пустым")
	}

	msg, err := ctx.Bot().Send(ctx.Chat(), "🔗 Получен токен. Авторизация…")
	if err != nil {
		return err
	}

	h.cache.RemoveLastAction(context.Background(), user.TelegramID) //nolint:errcheck
	bot, err := h.service.ConnectBot(context.Background(), &service.ConnectBotInput{
		UserID:   user.ID,
		BotToken: token,
	})
	if err != nil {
		if errors.Is(err, service.ErrBotAlreadyExists) {
			_, err := ctx.Bot().Edit(msg, "Этот бот уже добавлен в ваш аккаунт")
			if err != nil {
				return err
			}

			return nil
		}

		if errors.Is(err, service.ErrInvalidBot) {
			h.log.Warn(err.Error())
			_, err := ctx.Bot().Edit(msg, h.buildInvalidBotMessage())
			if err != nil {
				return err
			}

			return nil
		}

		return err
	}

	_, err = ctx.Bot().Edit(msg, "✅ Подключили бота, высылаем в него информацию")
	if err != nil {
		return err
	}

	_, err = bot.Send(&telebot.Chat{ID: user.TelegramID.Value()}, h.connectedBotWelcomeMessage())
	if err != nil {
		return err
	}

	message, markup, err := h.buildMainMenu(user)
	if err != nil {
		return err
	}

	return ctx.Send(message, markup)
}

func (h *Handler) connectedBotWelcomeMessage() string {
	hasModules := false

	if !hasModules {
		return `🤖 Ваш бот подключен к платформе — <b>Gargantua Labs</b>. На данный момент у вас не приобретены модули. Чтобы подключить, перейдите назад в бота платформы: /menu -> Баланс -> Выбрать модули`
	}

	return `
	🤖 Ваш бот подключен к платформе — <b>Gargantua Labs</b>. На данный момент у вас приобретены модули: [список оплаченных модулей]. Для настройки модулей, под ваши задачи, пропишите /start
	`
}

func (h *Handler) MyBotsCallback(ctx telebot.Context) error {
	if ctx.Callback() != nil {
		_ = ctx.Respond() //nolint:errcheck
	}

	user, ok := ctx.Get(middlewares.CtxUserKey).(*authDomain.User)
	if !ok {
		return ctx.Send("Не удалось определить пользователя.")
	}

	page := 1
	if ctx.Callback() != nil && ctx.Callback().Data != "" {
		if parsed, err := strconv.Atoi(ctx.Callback().Data); err == nil && parsed > 0 {
			page = parsed
		}
	}

	pagination, err := types.NewPagination(botsLimit, page)
	if err != nil {
		return shared.NewAppError(shared.CodeValidation, "failed parse pagination", err)
	}

	bots, totalCount, err := h.service.GetBotsByUserID(context.Background(), user.ID, &pagination)
	if err != nil {
		return err
	}

	if totalCount == 0 {
		markup := &telebot.ReplyMarkup{}
		connectBtn := markup.Data("➕ Подключить бота", string(constants.ConnectBotCallback))
		backBtn := markup.Data("⬅️ В меню", string(constants.ShowMainMenuCallback))
		markup.Inline(
			markup.Row(connectBtn),
			markup.Row(backBtn),
		)

		return ctx.Send("🤖 <b>У вас нет подключённых ботов</b>", &telebot.SendOptions{
			ParseMode:   telebot.ModeHTML,
			ReplyMarkup: markup,
		})
	}

	currentPage := pagination.Page.Value()
	totalPages := (totalCount + pagination.Limit.Value() - 1) / pagination.Limit.Value()

	// Формируем текст
	var sb strings.Builder
	sb.WriteString("🤖 <b>Ваши боты</b>\n\n")

	markup := &telebot.ReplyMarkup{}
	var rows []telebot.Row

	for i, b := range bots {
		botNumber := pagination.Limit.Value()*(currentPage-1) + i + 1

		sb.WriteString(fmt.Sprintf(
			"• <a href=\"https://t.me/%s\">%s</a> [#%d]\n",
			b.Name.Value(), b.Name.Value(), botNumber,
		))

		btn := markup.Data(
			fmt.Sprintf("#%d %s", botNumber, b.Name.Value()),
			string(constants.BotInfoCallback),
			b.ID.String(),
		)
		rows = append(rows, markup.Row(btn))
	}

	sb.WriteString(fmt.Sprintf("\n📄 Страница %d/%d", currentPage, totalPages))

	var navBtns []telebot.Btn
	if currentPage > 1 {
		prevBtn := markup.Data(
			"◀️ Назад",
			string(constants.MyBotsCallback),
			strconv.Itoa(currentPage-1),
		)
		navBtns = append(navBtns, prevBtn)
	}
	if currentPage < totalPages {
		nextBtn := markup.Data(
			"Вперёд ▶️",
			string(constants.MyBotsCallback),
			strconv.Itoa(currentPage+1),
		)
		navBtns = append(navBtns, nextBtn)
	}
	if len(navBtns) > 0 {
		rows = append(rows, markup.Row(navBtns...))
	}

	backBtn := markup.Data("⬅️ В меню", string(constants.ShowMainMenuCallback))
	rows = append(rows, markup.Row(backBtn))

	markup.Inline(rows...)

	opts := &telebot.SendOptions{
		ParseMode:             telebot.ModeHTML,
		ReplyMarkup:           markup,
		DisableWebPagePreview: true,
	}

	if ctx.Callback() != nil && ctx.Callback().Data != "" {
		return ctx.Edit(sb.String(), opts)
	}

	return ctx.Send(sb.String(), opts)
}

func (h *Handler) BotInfoCallback(ctx telebot.Context) error {
	if ctx.Callback() != nil {
		_ = ctx.Respond() //nolint:errcheck
	}

	bgCtx := context.Background()

	user, ok := ctx.Get(middlewares.CtxUserKey).(*authDomain.User)
	if !ok {
		return ctx.Send("Не удалось определить пользователя.")
	}

	botID, err := d.NewBotID(ctx.Callback().Data)
	if err != nil {
		return ctx.Send("Некорректный id бота")
	}

	bot, err := h.service.GetBotByUserID(bgCtx, user.ID, botID)
	if err != nil {
		return err
	}

	duration := formatDurationFull(bot.CreatedAt)
	activeText := "Неактивен"
	if bot.IsActive {
		activeText = "Активен"
	}

	markup := &telebot.ReplyMarkup{}
	backBtn := markup.Data("⬅️ К ботам", string(constants.MyBotsCallback), "1")
	menuBtn := markup.Data("⬅️ В меню", string(constants.ShowMainMenuCallback))
	markup.Inline(
		markup.Row(backBtn),
		markup.Row(menuBtn),
	)

	msg := fmt.Sprintf(`
🤖 <b>Информация о боте</b>

• Название: <a href="https://t.me/%s">%s</a>
• Дата подключения: %s (%s)
• Статус: [%s]
• Назначенные модули: [в разработке]
• RPS: мин 1, ср 1, макс 1
• Общая нагрузка: 0%%
• Кол-во пользователей: Общее: 1 / Активных: 1
`, bot.Name, bot.Name, bot.CreatedAt.Format("02.01.2006"), duration, activeText)

	return ctx.Send(msg, &telebot.SendOptions{
		DisableWebPagePreview: true,
		ParseMode:             telebot.ModeHTML,
		ReplyMarkup:           markup,
	})
}

func formatDurationFull(t time.Time) string {
	duration := time.Since(t)

	years := int(duration / (365 * 24 * time.Hour))
	duration -= time.Duration(years) * 365 * 24 * time.Hour

	days := int(duration / (24 * time.Hour))
	duration -= time.Duration(days) * 24 * time.Hour

	hours := int(duration / time.Hour)
	duration -= time.Duration(hours) * time.Hour

	minutes := int(duration / time.Minute)
	duration -= time.Duration(minutes) * time.Minute

	seconds := int(duration / time.Second)

	var parts []string

	if years > 0 {
		parts = append(parts, fmt.Sprintf("%dг", years))
	}
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%dд", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dч", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%dмин", minutes))
	}
	if seconds > 0 {
		parts = append(parts, fmt.Sprintf("%dс", seconds))
	}

	if len(parts) == 0 {
		return "0с"
	}

	return strings.Join(parts, ",")
}

func formatRps(f float64) string {
	if f == math.Trunc(f) {
		return fmt.Sprintf("%.0f", f)
	}
	return fmt.Sprintf("%.1f", f)
}
