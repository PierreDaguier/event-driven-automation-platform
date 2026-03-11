package repo

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisIdempotency struct {
	client *redis.Client
}

func NewRedisIdempotency(addr, password string) *RedisIdempotency {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
	return &RedisIdempotency{client: client}
}

func (r *RedisIdempotency) Reserve(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	ok, err := r.client.SetNX(ctx, key, "1", ttl).Result()
	return ok, err
}

func (r *RedisIdempotency) Close() error {
	return r.client.Close()
}
