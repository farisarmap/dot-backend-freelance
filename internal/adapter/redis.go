package adapter

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheManager interface {
	Get(key string) (string, error)
	Set(key string, value interface{}) error
	Delete(key string) error
}

type redisCache struct {
	client  *redis.Client
	ctx     context.Context
	timeout time.Duration
}

func NewRedisCache(client *redis.Client, ttl time.Duration) CacheManager {
	return &redisCache{
		client:  client,
		ctx:     context.Background(),
		timeout: ttl,
	}
}

func (r *redisCache) Get(key string) (string, error) {
	val, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (r *redisCache) Set(key string, value interface{}) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(r.ctx, key, bytes, r.timeout).Err()
}

func (r *redisCache) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}
