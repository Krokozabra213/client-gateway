package config

import "time"

type HTTPConfig struct {
	Host           string        `yaml:"host" env:"HTTP_HOST" env-default:"0.0.0.0"`
	Port           string        `yaml:"port" env:"HTTP_PORT" env-default:"8080"`
	ReadTimeout    time.Duration `yaml:"readTimeout" env:"HTTP_READ_TIMEOUT" env-default:"10s"`
	WriteTimeout   time.Duration `yaml:"writeTimeout" env:"HTTP_WRITE_TIMEOUT" env-default:"10s"`
	MaxHeaderBytes int           `yaml:"maxHeaderBytes" env:"HTTP_MAX_HEADER_BYTES" env-default:"1"`
	BodyLimitBytes int           `yaml:"bodyLimitBytes" env:"HTTP_BODY_LIMIT_BYTES" env-default:"4"`
	AllowedOrigins []string      `yaml:"allowedOrigins" env:"ALLOWED_ORIGINS" env-separator:","`
}

func (c HTTPConfig) HTTPAddress() string {
	return c.Host + ":" + c.Port
}
