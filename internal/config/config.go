package config

import (
	"fmt"
	"log"
	"os"

	"github.com/Krokozabra213/gargantua_common/pkg/config"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	App          AppConfig             `yaml:"app"`
	Postgres     config.PostgresConfig `yaml:"postgres"`
	HTTP         config.HTTPConfig     `yaml:"http"`
	Bot          BotConfig             `yaml:"bot"`
	Redis        RedisConfig           `yaml:"redis"`
	Logger       config.SlogConfig     `yaml:"slog"`
	Payment      PaymentConfig         `yaml:"payment"`
	Balance      BalanceConfig         `yaml:"balance"`
	OrderProduct OrderProductConfig    `yaml:"order_product"`
}

func Init() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Print(err)
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}

	configFile := fmt.Sprintf("configs/%s.yaml", env)

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s (ENV=%s)", configFile, env)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configFile, &cfg); err != nil {
		return nil, err
	}

	err := cfg.Validate()
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) Validate() error {
	return c.Logger.Validate()
}
