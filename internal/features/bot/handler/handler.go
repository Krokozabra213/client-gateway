package handler

import (
	"context"
	"log/slog"
	"time"

	"github.com/GargantuaLabs/client-gateway/internal/clients"
	"github.com/GargantuaLabs/client-gateway/internal/config"
	authDomain "github.com/GargantuaLabs/client-gateway/internal/features/auth/domain"
	"github.com/GargantuaLabs/client-gateway/internal/features/auth/service"
	"github.com/GargantuaLabs/client-gateway/internal/features/bot/constants"
	d "github.com/GargantuaLabs/client-gateway/internal/features/bot/domain"
	"github.com/GargantuaLabs/client-gateway/internal/features/bot/middlewares"
	botService "github.com/GargantuaLabs/client-gateway/internal/features/bot/service"
	"github.com/GargantuaLabs/client-gateway/internal/features/shared/types"

	"gopkg.in/telebot.v4"
)

type AuthService interface {
	Register(ctx context.Context, data *service.RegisterInput) error
	GetUserByTelegramID(ctx context.Context, telegramID types.TelegramID) (*authDomain.User, error)
}

type Service interface {
	ConnectBot(ctx context.Context, botToken *botService.ConnectBotInput) (*telebot.Bot, error)
	GetBotsByUserID(ctx context.Context, userID authDomain.UserID, pagination *types.Pagination) ([]*d.Bot, int, error)
	GetConnectedBot(ctx context.Context, botID d.BotID) (*telebot.Bot, authDomain.UserID, error)
	GetBotByUserID(ctx context.Context, userID authDomain.UserID, botID d.BotID) (*d.Bot, error)
	GetBotCount(ctx context.Context, userID authDomain.UserID) (int, error)
	GetBotLastMinuteRps(ctx context.Context, botID d.BotID) (min, avg, max float64, err error)
}

type Cache interface {
	GetAndDeleteGender(ctx context.Context, telegramID types.TelegramID) (authDomain.UserGender, error)
	StoreGender(ctx context.Context, telegramID types.TelegramID, gender authDomain.UserGender) error

	StoreLastAction(ctx context.Context, lastAction constants.LastAction, senderID types.TelegramID) error
	GetLastAction(ctx context.Context, senderID types.TelegramID) (constants.LastAction, error)
	RemoveLastAction(ctx context.Context, senderID types.TelegramID) error

	StoreMessageID(ctx context.Context, senderID types.TelegramID, messageID int, ttl *time.Duration) error
	GetAndDeleteMessageID(ctx context.Context, senderID types.TelegramID) (int, error)
}

type Handler struct {
	authService  AuthService
	service      Service
	cache        Cache
	payment      *clients.PaymentClient
	balance      *clients.BalanceClient
	orderProduct *clients.OrderProductClient
	cfg          *config.BotConfig
	log          *slog.Logger
	middlewares  *middlewares.Middlewares
	markups      map[constants.TelegramMarkup]*telebot.ReplyMarkup
}

func New(
	bot *telebot.Bot, authService AuthService, service Service, cache Cache,
	cfg *config.BotConfig, log *slog.Logger, middlewares *middlewares.Middlewares, payment *clients.PaymentClient,
	balance *clients.BalanceClient, orderProduct *clients.OrderProductClient,
) *Handler {
	handler := &Handler{
		payment:      payment,
		authService:  authService,
		service:      service,
		cache:        cache,
		cfg:          cfg,
		log:          log,
		middlewares:  middlewares,
		markups:      make(map[constants.TelegramMarkup]*telebot.ReplyMarkup),
		balance:      balance,
		orderProduct: orderProduct,
	}

	handler.initMarkups(bot)
	return handler
}

func (h *Handler) initMarkups(bot *telebot.Bot) {
	registerSelector := &telebot.ReplyMarkup{}
	btnRegister := registerSelector.Data("РЕГИСТРАЦИЯ", string(constants.RegisterCallback))
	registerSelector.Inline(
		registerSelector.Row(btnRegister),
	)
	h.markups[constants.StartupMarkup] = registerSelector
	bot.Handle(&btnRegister, h.registerCallback)

	registerConfirmSelector := &telebot.ReplyMarkup{}
	btnRegisterApprove := registerConfirmSelector.Data("✅ Я согласен", string(constants.ApproveRegisterCallback))
	btnRegisterCancel := registerConfirmSelector.Data("Я несогласен 🚫", string(constants.CancelRegisterCallback))
	registerConfirmSelector.Inline(
		registerConfirmSelector.Row(btnRegisterApprove, btnRegisterCancel),
	)
	h.markups[constants.RegisterMarkup] = registerConfirmSelector
	bot.Handle(&btnRegisterApprove, h.registerApproveCallback)
	bot.Handle(&btnRegisterCancel, h.registerCancelCallback)

	genderSelector := &telebot.ReplyMarkup{}
	btnMale := genderSelector.Data("🧑 Мужчина", string(constants.MaleCallback))
	btnFemale := genderSelector.Data("Женщина 👩", string(constants.FemaleCallback))
	genderSelector.Inline(
		genderSelector.Row(btnMale, btnFemale),
	)
	h.markups[constants.GenderMarkup] = genderSelector
	bot.Handle(&btnMale, h.genderCallback)
	bot.Handle(&btnFemale, h.genderCallback)

	birthdaySelector := &telebot.ReplyMarkup{}

	buttons := [][][2]string{
		{{"1", "b_1"}, {"2", "b_2"}, {"3", "b_3"}},
		{{"4", "b_4"}, {"5", "b_5"}, {"6", "b_6"}},
		{{"7", "b_7"}, {"8", "b_8"}, {"9", "b_9"}},
		{{"0", "b_0"}},
		{{"✅", "b_ok"}, {"❌", "b_del"}},
	}

	rows := make([]telebot.Row, len(buttons))
	for i, row := range buttons {
		btns := make([]telebot.Btn, len(row))
		for j, b := range row {
			btn := birthdaySelector.Data(b[0], b[1], b[0])
			btns[j] = btn
			bot.Handle(&btn, h.birthdayCallback)
		}
		rows[i] = birthdaySelector.Row(btns...)
	}

	birthdaySelector.Inline(rows...)
	h.markups[constants.BirthdayMarkup] = birthdaySelector

	authGroup := bot.Group()
	authGroup.Use(h.middlewares.Authenticate)

	registerSuccessSelector := &telebot.ReplyMarkup{}
	btnMainMenu := registerSuccessSelector.Data("ГЛАВНОЕ МЕНЮ", string(constants.ShowMainMenuCallback))
	btnHelp := registerSuccessSelector.Data("ПОМОЩЬ", string(constants.HelpCallback))
	registerSuccessSelector.Inline(
		registerSuccessSelector.Row(btnMainMenu, btnHelp),
	)
	h.markups[constants.RegisterSuccessMarkup] = registerSuccessSelector
	authGroup.Handle(&btnMainMenu, h.mainMenuCallback)
	authGroup.Handle(&btnHelp, h.helpCallback)

	mainMenuSelector := &telebot.ReplyMarkup{}
	connectBotBtn := mainMenuSelector.Data("➕ Подключить бота", string(constants.ConnectBotCallback))
	managementBtn, statisticsBtn, balanceBtn, helpBtn := makeCommonButtons(mainMenuSelector)
	myOrdersBtn := mainMenuSelector.Data("📦 Заказы", string(constants.MyOrdersCallback))
	mySubscriptionsBtn := mainMenuSelector.Data("⭐ Подписки", string(constants.MySubscriptionsCallback))
	shopBtn := mainMenuSelector.Data("⭐ Магазин подписок", string(constants.ShopCallback))

	mainMenuSelector.Inline(
		mainMenuSelector.Row(connectBotBtn),
		mainMenuSelector.Row(managementBtn, statisticsBtn),
		mainMenuSelector.Row(myOrdersBtn, mySubscriptionsBtn),
		mainMenuSelector.Row(shopBtn),
		mainMenuSelector.Row(balanceBtn),
		mainMenuSelector.Row(helpBtn),
	)
	h.markups[constants.MainMenuMarkup] = mainMenuSelector
	authGroup.Handle(&connectBotBtn, h.connectBotCallback)
	authGroup.Handle(&managementBtn, h.connectBotCallback)
	authGroup.Handle(&statisticsBtn, h.connectBotCallback)
	authGroup.Handle(&balanceBtn, h.balanceCallback)
	authGroup.Handle(&helpBtn, h.helpCallback)

	mainMenuWithBots := &telebot.ReplyMarkup{}
	connectBtn := mainMenuWithBots.Data("➕ Подключить", string(constants.ConnectBotCallback))
	myBotsBtn := mainMenuWithBots.Data("🤖 Мои боты", string(constants.MyBotsCallback))
	myOrdersBtn2 := mainMenuWithBots.Data("📦 Заказы", string(constants.MyOrdersCallback))
	mySubscriptionsBtn2 := mainMenuWithBots.Data("⭐ Подписки", string(constants.MySubscriptionsCallback))
	shopBtn2 := mainMenuWithBots.Data("⭐ Магазин подписок", string(constants.ShopCallback))
	managementBtn2, statisticsBtn2, balanceBtn2, helpBtn2 := makeCommonButtons(mainMenuSelector)

	mainMenuWithBots.Inline(
		mainMenuWithBots.Row(connectBtn, myBotsBtn),
		mainMenuWithBots.Row(managementBtn2, statisticsBtn2),
		mainMenuWithBots.Row(myOrdersBtn2, mySubscriptionsBtn2),
		mainMenuWithBots.Row(shopBtn2),
		mainMenuWithBots.Row(balanceBtn2),
		mainMenuWithBots.Row(helpBtn2),
	)
	h.markups[constants.MainMenuWithBotsMarkup] = mainMenuWithBots
	authGroup.Handle(&connectBtn, h.connectBotCallback)
	//authGroup.Handle(&myBotsBtn, h.myBotsCallback(authGroup))
	authGroup.Handle(&myBotsBtn, h.MyBotsCallback)
	authGroup.Handle(&shopBtn, h.ShopCallback)
	authGroup.Handle(&helpBtn2, h.helpCallback)
	authGroup.Handle(&managementBtn2, h.managementCallback)
	authGroup.Handle(&statisticsBtn2, h.statisticsCallback)

	connectBotSelector := &telebot.ReplyMarkup{}
	connectBotConfirmBtn := connectBotSelector.Data("Приступить", string(constants.ConnectBotConfirmCallback))
	backToMenuBtn := connectBotSelector.Data("⬅️ Назад", string(constants.ShowMainMenuCallback))
	connectBotSelector.Inline(
		connectBotSelector.Row(connectBotConfirmBtn),
		connectBotSelector.Row(backToMenuBtn),
	)
	h.markups[constants.ConnectBotConfirmMarkup] = connectBotSelector
	authGroup.Handle(&connectBotConfirmBtn, h.connectBotConfirmCallback)

	shopSelector := &telebot.ReplyMarkup{}
	backBtn := shopSelector.Data("⬅️ Назад", string(constants.ShowMainMenuCallback))
	shopSelector.Inline(shopSelector.Row(backBtn))
	h.markups[constants.ShopMarkup] = shopSelector

	balanceSelector := &telebot.ReplyMarkup{}

	depositBtn := balanceSelector.Data("💰 Пополнить баланс", string(constants.DepositCallback))
	chooseModuleBtn := balanceSelector.Data("⚙️ Выбрать модули", string(constants.ManagementCallback))
	balanceSelector.Inline(
		balanceSelector.Row(depositBtn),
		balanceSelector.Row(chooseModuleBtn),
		balanceSelector.Row(backToMenuBtn),
	)
	h.markups[constants.BalanceMarkup] = balanceSelector

	authGroup.Handle(&depositBtn, h.depositCallback)

	depositSelector := &telebot.ReplyMarkup{}
	tgDepositBtn := depositSelector.Data("Пополнить звездами TG", string(constants.TGDepositCallback))
	tbankDepositBtn := depositSelector.Data("Пополнить через TBank", string(constants.TbankDepositCallback))
	backToMenuBtn2 := depositSelector.Data("⬅️ В меню", string(constants.ShowMainMenuCallback))
	depositSelector.Inline(
		depositSelector.Row(tgDepositBtn),
		depositSelector.Row(tbankDepositBtn),
		depositSelector.Row(backToMenuBtn2),
	)
	h.markups[constants.DepositMarkup] = depositSelector

	authGroup.Handle(&tgDepositBtn, h.TGDepositCallback)
	authGroup.Handle(&tbankDepositBtn, h.tbankDepositCallback)
}

func makeCommonButtons(m *telebot.ReplyMarkup) (management, stats, balance, help telebot.Btn) {
	management = m.Data("⚙️ Управление", string(constants.ManagementCallback))
	stats = m.Data("📊 Статистика", string(constants.StatisticsCallback))
	balance = m.Data("💰 Баланс", string(constants.BalanceCallback))
	help = m.Data("❓ Помощь", string(constants.HelpCallback))

	return
}
