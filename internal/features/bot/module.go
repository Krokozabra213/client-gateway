package bot

import (
	"log/slog"

	"github.com/GargantuaLabs/client-gateway/internal/clients"
	"github.com/GargantuaLabs/client-gateway/internal/config"
	"github.com/GargantuaLabs/client-gateway/internal/features/bot/constants"
	errorhandler "github.com/GargantuaLabs/client-gateway/internal/features/bot/error_handler"
	"github.com/GargantuaLabs/client-gateway/internal/features/bot/handler"
	"github.com/GargantuaLabs/client-gateway/internal/features/bot/middlewares"
	postgresRepo "github.com/GargantuaLabs/client-gateway/internal/features/bot/repository/postgres"
	redisRepo "github.com/GargantuaLabs/client-gateway/internal/features/bot/repository/redis"
	"github.com/GargantuaLabs/client-gateway/internal/features/bot/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"

	"github.com/gofiber/fiber/v3"
	"gopkg.in/telebot.v4"
	"gopkg.in/telebot.v4/middleware"
)

type Module struct {
	handler     *handler.Handler
	bot         *telebot.Bot
	botCfg      *config.BotConfig
	log         *slog.Logger
	middlewares *middlewares.Middlewares
}

type Dependencies struct {
	Payment           *clients.PaymentClient
	Balance           *clients.BalanceClient
	OrderProduct      *clients.OrderProductClient
	BotClient         *telebot.Bot
	BotCfg            *config.BotConfig
	Redis             *redis.Client
	AuthService       service.AuthService
	EncryptionService service.EncryptionService
	HashingService    service.HashingService
	Log               *slog.Logger
	PGXPool           *pgxpool.Pool
	Registerer        prometheus.Registerer
}

func New(deps *Dependencies) *Module {
	postgresR := postgresRepo.New(deps.PGXPool)
	redisR := redisRepo.New(deps.Redis)
	s := service.New(deps.AuthService, deps.EncryptionService, deps.HashingService, postgresR, redisR, deps.BotCfg)
	m := middlewares.New(deps.AuthService, s, deps.Registerer, deps.Log)
	h := handler.New(deps.BotClient, deps.AuthService, s, redisR, deps.BotCfg, deps.Log, m,
		deps.Payment, deps.Balance, deps.OrderProduct)

	return &Module{
		handler:     h,
		bot:         deps.BotClient,
		botCfg:      deps.BotCfg,
		middlewares: m,
		log:         deps.Log,
	}
}

func (m *Module) HTTPv1() func(fiber.Router) {
	errHandler := errorhandler.NewBotErrorHandler(m.log)
	m.bot.Use(middleware.Recover(func(err error, ctx telebot.Context) {
		errHandler.Handle(err, ctx)
	}))

	m.bot.Handle("/start", m.handler.Start)

	m.bot.Use(m.middlewares.Authenticate)
	m.bot.Handle("/menu", m.handler.Menu)
	m.bot.Handle(telebot.OnText, m.handler.HandleMessage)
	m.bot.Handle(telebot.OnCheckout, m.handler.HandleTGCheckout)
	m.bot.Handle(telebot.OnPayment, m.handler.HandleTGSuccessful)
	m.bot.Handle("\f"+string(constants.ShopItemCallback), m.handler.ShopProductCallback)
	m.bot.Handle("\f"+string(constants.BuyProductCallback), m.handler.BuyProductCallback)
	m.bot.Handle("\f"+string(constants.MyOrdersCallback), m.handler.OrdersCallback)
	m.bot.Handle("\f"+string(constants.MySubscriptionsCallback), m.handler.SubscriptionsCallback)
	m.bot.Handle("\f"+string(constants.MyBotsCallback), m.handler.MyBotsCallback)
	m.bot.Handle("\f"+string(constants.BotInfoCallback), m.handler.BotInfoCallback)
	m.bot.Handle("\f"+string(constants.TGDepositAmountCallback), m.handler.TGDepositAmountCallback)
	m.bot.Handle("\f"+string(constants.TbankDepositAmountCallback), m.handler.TbankDepositAmountCallback)

	return func(router fiber.Router) {
		router.Post("/", m.handler.ResponseInternalWebhook(m.bot, m.botCfg))
		router.Post("/:botId", m.handler.RespondExternalWebhook(m.middlewares, m.botCfg))
	}
}
