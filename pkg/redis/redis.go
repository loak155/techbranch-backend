package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisManager struct {
	expiration time.Duration
	client     *redis.Client
}

func NewRedisManager(addr string, db int, expiration time.Duration) *RedisManager {
	return &RedisManager{
		expiration: expiration,
		client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: "",
			DB:       db,
		}),
	}
}

func (rm *RedisManager) Set(ctx context.Context, key string, value interface{}) error {
	err := rm.client.Set(ctx, key, value, rm.expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set key: %v", err)
	}
	return nil
}

func (rm *RedisManager) Get(ctx context.Context, key string) (string, error) {
	val, err := rm.client.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("failed to get key: %v", err)
	}
	return val, nil
}

func (rm *RedisManager) Del(ctx context.Context, key string) error {
	err := rm.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete key: %v", err)
	}
	return nil
}
