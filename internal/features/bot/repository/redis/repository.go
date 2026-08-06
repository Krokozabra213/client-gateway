package redisRepo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	authDomain "github.com/GargantuaLabs/client-gateway/internal/features/auth/domain"
	"github.com/GargantuaLabs/client-gateway/internal/features/bot/constants"
	"github.com/GargantuaLabs/client-gateway/internal/features/bot/domain"
	"github.com/GargantuaLabs/client-gateway/internal/features/shared/types"
	"github.com/google/uuid"

	"github.com/redis/go-redis/v9"
)

type botCacheData struct {
	Token   types.NonEmptyString `json:"token"`
	OwnerID authDomain.UserID    `json:"owner_id"`
}

type Cache struct {
	redis *redis.Client
}

func (c *Cache) GetAndDeleteGender(ctx context.Context, telegramID types.TelegramID) (authDomain.UserGender, error) {
	gender, err := c.redis.GetDel(ctx, generateGenderKey(telegramID)).Result()
	if err != nil {
		return "", err
	}

	parsed, err := authDomain.NewUserGender(gender)
	if err != nil {
		return "", err
	}

	return parsed, nil
}

func (c *Cache) StoreGender(ctx context.Context, telegramID types.TelegramID, gender authDomain.UserGender) error {
	return c.redis.Set(ctx, generateGenderKey(telegramID), string(gender), time.Minute*5).Err()
}

func (c *Cache) StoreLastAction(ctx context.Context, lastAction constants.LastAction, senderID types.TelegramID) error {
	return c.redis.Set(ctx, generateLastActionKey(senderID), string(lastAction), time.Minute*5).Err()
}

func (c *Cache) GetLastAction(ctx context.Context, senderID types.TelegramID) (constants.LastAction, error) {
	result, err := c.redis.Get(ctx, generateLastActionKey(senderID)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", err
	}

	return constants.LastAction(result), err
}

func (c *Cache) RemoveLastAction(ctx context.Context, senderID types.TelegramID) error {
	return c.redis.Del(ctx, generateLastActionKey(senderID)).Err()
}

func (c *Cache) StoreMessageID(ctx context.Context, senderID types.TelegramID, messageID int, ttl *time.Duration) error {
	var t time.Duration

	if ttl == nil {
		t = time.Minute * 5
	} else {
		t = *ttl
	}

	return c.redis.Set(ctx, generateMessageKey(senderID), messageID, t).Err()
}

func (c *Cache) GetAndDeleteMessageID(ctx context.Context, senderID types.TelegramID) (int, error) {
	val, err := c.redis.GetDel(ctx, generateMessageKey(senderID)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}

		return 0, err
	}

	id, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (c *Cache) StoreBotCount(ctx context.Context, count int, ownerID authDomain.UserID) error {
	return c.redis.Set(ctx, generateBotCountKey(ownerID), count, time.Minute*60).Err()
}

func (c *Cache) GetBotCount(ctx context.Context, ownerID authDomain.UserID) (int, error) {
	val, err := c.redis.Get(ctx, generateBotCountKey(ownerID)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return 0, err
	}

	if val == "" {
		return 0, nil
	}

	converted, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	return converted, nil
}

func (c *Cache) IncrementBotCount(ctx context.Context, ownerID authDomain.UserID) error {
	return c.redis.Incr(ctx, generateBotCountKey(ownerID)).Err()
}

func (c *Cache) StoreBotData(ctx context.Context, botID domain.BotID, botToken types.NonEmptyString, ownerID authDomain.UserID) error {
	data := botCacheData{
		Token:   botToken,
		OwnerID: ownerID,
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return c.redis.Set(ctx, generateBotDataKey(botID), bytes, time.Minute*10).Err()
}

func (c *Cache) GetBotData(ctx context.Context, botID domain.BotID) (botToken types.NonEmptyString, ownerID authDomain.UserID, err error) {
	val, err := c.redis.Get(ctx, generateBotDataKey(botID)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return types.NonEmptyString{}, authDomain.UserID(uuid.Nil), nil
		}
		return types.NonEmptyString{}, authDomain.UserID(uuid.Nil), err
	}

	var data botCacheData
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return types.NonEmptyString{}, authDomain.UserID(uuid.Nil), err
	}

	return data.Token, data.OwnerID, nil
}

func (c *Cache) CreateBotRps(ctx context.Context, botID domain.BotID) error {
	_, err := c.redis.Do(ctx, "TS.CREATE", generateBotRpsKey(botID),
		"RETENTION", 86400000,
		"LABELS", "type", "rps",
	).Result()

	if err != nil && !strings.Contains(err.Error(), "already exists") {
		return err
	}

	return nil
}

func (c *Cache) IncrementBotRps(ctx context.Context, botID domain.BotID) error {
	key := generateBotRpsKey(botID)
	pipe := c.redis.Pipeline()

	pipe.Do(ctx, "TS.ADD", key, "*", 1, "RETENTION", 120_000)
	pipe.Expire(ctx, key, 2*time.Minute)

	_, err := pipe.Exec(ctx)

	return err
}

func (c *Cache) GetBotLastMinuteRps(ctx context.Context, botID domain.BotID) (min, avg, max float64, err error) {
	now := time.Now().UnixMilli()
	oneMinuteAgo := now - 60_000

	res, err := c.redis.Do(ctx, "TS.RANGE", generateBotRpsKey(botID),
		oneMinuteAgo, now,
		"AGGREGATION", "sum", 1000,
	).Result()
	if err != nil {
		return 0, 0, 0, err
	}

	points, ok := res.([]interface{})
	if !ok || len(points) == 0 {
		return 0, 0, 0, nil
	}

	min = math.MaxFloat64
	max = -math.MaxFloat64
	var sum float64 = 0

	for _, p := range points {
		point, ok := p.([]interface{})
		if !ok || len(point) < 2 {
			continue
		}

		val, errConv := strconv.ParseFloat(fmt.Sprintf("%v", point[1]), 64)
		if errConv != nil {
			continue
		}

		if val < min {
			min = val
		}

		if val > max {
			max = val
		}
		sum += val
	}

	avg = sum / float64(len(points))
	return min, avg, max, nil
}

func generateGenderKey(senderID types.TelegramID) string {
	return fmt.Sprintf("telegram_gender:%d", senderID.Value())
}
func generateLastActionKey(senderID types.TelegramID) string {
	return fmt.Sprintf("telegram_last_action:%d", senderID.Value())
}
func generateMessageKey(senderID types.TelegramID) string {
	return fmt.Sprintf("telegram_message:%d", senderID.Value())
}
func generateBotCountKey(ownerID authDomain.UserID) string {
	return fmt.Sprintf("bot_count:%s", ownerID)
}
func generateBotDataKey(botID domain.BotID) string {
	return fmt.Sprintf("bot_data:%s", botID)
}

func generateBotRpsKey(botID domain.BotID) string {
	return fmt.Sprintf("bot_rps:%s", botID)
}

func New(redis *redis.Client) *Cache {
	return &Cache{redis}
}
