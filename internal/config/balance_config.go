package config

type BalanceConfig struct {
	Addr     string `yaml:"addr" env-required:"true"`
}
