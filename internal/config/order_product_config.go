package config

type OrderProductConfig struct {
	Addr string `yaml:"addr" env-required:"true"`
}
