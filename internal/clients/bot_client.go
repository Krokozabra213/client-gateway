package clients

import (
	"log/slog"

	"github.com/GargantuaLabs/client-gateway/internal/config"
	errorhandler "github.com/GargantuaLabs/client-gateway/internal/features/bot/error_handler"
	"gopkg.in/telebot.v4"
)

func NewBotClient(cfg *config.BotConfig, log *slog.Logger) (*telebot.Bot, error) {
	var b *telebot.Bot
	var err error

	if log == nil {
		b, err = telebot.NewBot(telebot.Settings{
			Token:     cfg.Token,
			ParseMode: telebot.ModeHTML,
		})
	} else {
		errHandler := errorhandler.NewBotErrorHandler(log)
		b, err = telebot.NewBot(telebot.Settings{
			Token:     cfg.Token,
			ParseMode: telebot.ModeHTML,
			OnError:   errHandler.Handle,
		})
	}

	if err != nil {
		return nil, err
	}

	if err := b.RemoveWebhook(true); err != nil {
		return nil, err
	}

	if err := b.SetWebhook(&telebot.Webhook{
		Endpoint:    &telebot.WebhookEndpoint{PublicURL: cfg.WebhookUrl},
		SecretToken: cfg.Secret,
	}); err != nil {
		return nil, err
	}

	return b, nil
}
