package config

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

type PostgresConfig struct {
	Host     string `yaml:"host" env:"POSTGRES_HOST" env-required:"true"`
	Port     string `yaml:"port" env:"POSTGRES_PORT" env-required:"true"`
	User     string `yaml:"user" env:"POSTGRES_USER" env-required:"true"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD" env-required:"true"`
	DBName   string `yaml:"dbName" env:"POSTGRES_DB" env-required:"true"`

	SSLMode        string        `yaml:"sslMode" env:"PG_SSL_MODE" env-default:"disable"`
	ConnectTimeout time.Duration `yaml:"connectTimeout" env:"PG_CONNECT_TIMEOUT" env-default:"5s"`
	MaxConns       int           `yaml:"maxConns" env:"PG_MAX_CONNS" env-default:"20"`
	MinConns       int           `yaml:"minConns" env:"PG_MIN_CONNS" env-default:"2"`
}

func (c PostgresConfig) PGAddress() string {
	return c.Host + ":" + c.Port
}

func (c *PostgresConfig) DSN() string {
	user := url.PathEscape(c.User)
	password := url.PathEscape(c.Password)
	dbName := url.PathEscape(c.DBName)

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		user, password, c.Host, c.Port, dbName)

	var params []string
	if c.SSLMode != "" {
		params = append(params, "sslmode="+c.SSLMode)
	}

	if c.ConnectTimeout > 0 {
		params = append(params, fmt.Sprintf("connect_timeout=%d", int(c.ConnectTimeout.Seconds())))
	}

	if len(params) > 0 {
		dsn += "?" + strings.Join(params, "&")
	}

	return dsn
}
