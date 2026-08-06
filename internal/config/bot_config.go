package config

type BotConfig struct {
	Token      string `yaml:"token" env:"BOT_TOKEN" env-required:"true"`
	WebhookUrl string `yaml:"webhookUrl" env:"BOT_WEBHOOK_URL" env-required:"true"`
	Secret     string `yaml:"secret" env:"BOT_SECRET" env-required:"true"`
}
