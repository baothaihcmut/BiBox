package lock

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrKeyNotFound   = errors.New("lock key not found")
	ErrValueMisMatch = errors.New("lock value mismatch")
)

type RedisDistributedLockService struct {
	client *redis.Client
}

func (r *RedisDistributedLockService) AcquireLock(ctx context.Context, key string, value string, maxRetry int, interval time.Duration, ttl time.Duration) (bool, error) {
	for range maxRetry {
		<-time.After(interval)
		ok, err := r.client.SetNX(ctx, key, value, ttl).Result()
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}

func (r *RedisDistributedLockService) ReleaseLock(ctx context.Context, key string, value string) error {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return ErrKeyNotFound
		}
		return err
	}
	if val != value {
		return ErrValueMisMatch
	}
	if _, err := r.client.Del(ctx, key).Result(); err != nil {
		return err
	}
	return nil
}
func NewRedisDistributedLockService(c *redis.Client) *RedisDistributedLockService {
	return &RedisDistributedLockService{
		client: c,
	}
}
