package clients

import (
	"context"
	"fmt"

	"github.com/GargantuaLabs/client-gateway/internal/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg *config.RedisConfig) (*redis.Client, error) {
	options := &redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.Database,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	}

	rdb := redis.NewClient(options)
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return rdb, nil
}
