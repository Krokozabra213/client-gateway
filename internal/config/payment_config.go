package config

type PaymentConfig struct {
	Addr     string `yaml:"addr" env-required:"true"`
}
