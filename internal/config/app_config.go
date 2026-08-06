package config

type AppConfig struct {
	Name                string `yaml:"name" env-default:"gateway-api"`
	UserCacheTTLMinutes int    `yaml:"userCacheTTLMinutes" env:"USER_CACHE_TTL_MINUTES" env-required:"true"`
	EncryptionKey       string `yaml:"encryptionKey" env:"ENCRYPTION_KEY" env-required:"true"`
	HashingKey          string `yaml:"hashingKey" env:"HASHING_KEY" env-required:"true"`
}
