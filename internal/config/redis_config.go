package config

type RedisConfig struct {
	Addr     string `yaml:"addr" env:"REDIS_ADDR" env-required:"true"`
	Password string `yaml:"password" env:"REDIS_PASSWORD" env-required:"true"`
	Database int    `yaml:"database" env:"REDIS_DATABASE" env-default:"0"`

	PoolSize     int `yaml:"poolSize" env:"REDIS_POOL_SIZE" env-default:"10"`
	MinIdleConns int `yaml:"minIdleConns" env:"REDIS_MIN_IDLE_CONNS" env-default:"2"`
}
