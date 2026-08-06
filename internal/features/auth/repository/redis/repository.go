package redisRepo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	d "github.com/GargantuaLabs/client-gateway/internal/features/auth/domain"
	"github.com/GargantuaLabs/client-gateway/internal/features/shared/types"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	redis *redis.Client
}

func (r *Cache) GetUser(ctx context.Context, telegramID types.TelegramID) (*d.User, error) {
	cached, err := r.redis.Get(ctx, generateUserKey(telegramID)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if cached == "" {
		return nil, nil
	}

	var user d.User
	if err := json.Unmarshal([]byte(cached), &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Cache) SetUser(ctx context.Context, user *d.User, ttl time.Duration) error {
	str, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return r.redis.Set(ctx, generateUserKey(user.TelegramID), string(str), ttl).Err()
}

func generateUserKey(telegramID types.TelegramID) string {
	return fmt.Sprintf("user:%d", telegramID.Value())
}

func New(redis *redis.Client) *Cache {
	return &Cache{redis}
}
