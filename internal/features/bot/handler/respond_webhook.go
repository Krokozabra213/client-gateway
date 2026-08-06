package handler

import (
	"context"
	"net/http"

	"github.com/GargantuaLabs/client-gateway/internal/config"
	d "github.com/GargantuaLabs/client-gateway/internal/features/bot/domain"
	"github.com/GargantuaLabs/client-gateway/internal/features/bot/middlewares"
	"github.com/gofiber/fiber/v3"
	"gopkg.in/telebot.v4"
)

func (h *Handler) ResponseInternalWebhook(bot *telebot.Bot, botCfg *config.BotConfig) func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		var update telebot.Update

		secret := c.Get("X-Telegram-Bot-Api-Secret-Token")
		if secret != botCfg.Secret {
			return c.SendStatus(http.StatusOK)
		}

		if err := c.Bind().All(&update); err != nil {
			return c.SendStatus(http.StatusOK)
		}

		if update.Callback != nil {
			data := update.Callback.Data
			h.log.Debug("pressed button...", "button_data", data)
		}

		bot.ProcessUpdate(update)
		return c.SendStatus(http.StatusOK)
	}
}

func (h *Handler) RespondExternalWebhook(m *middlewares.Middlewares, botCfg *config.BotConfig) func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		id := c.Params("botId")
		botID, err := d.NewBotID(id)
		if err != nil {
			return c.SendStatus(http.StatusOK)
		}

		var update telebot.Update

		secret := c.Get("X-Telegram-Bot-Api-Secret-Token")
		if secret != botCfg.Secret {
			return c.SendStatus(http.StatusOK)
		}

		if err := c.Bind().All(&update); err != nil {
			return c.SendStatus(http.StatusOK)
		}

		bot, userID, err := h.service.GetConnectedBot(context.Background(), botID)
		if err != nil {
			return c.SendStatus(http.StatusOK)
		}

		bot.Use(m.Authenticate, m.OwnerOnlyMiddleware(userID), m.CalculateBotRps(botID))
		bot.Handle("/start", func(ctx telebot.Context) error {
			return ctx.Send("START")
		})
		bot.ProcessUpdate(update)
		return c.SendStatus(http.StatusOK)
	}
}
