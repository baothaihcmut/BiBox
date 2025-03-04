package cache

import (
	"context"
	"time"
)

type CacheService interface {
	SetValue(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	GetValue(ctx context.Context, key string, output interface{}) error
	SetString(ctx context.Context, key string, val string, ttl time.Duration) error
	GetString(ctx context.Context, key string) (*string, error)
	Remove(ctx context.Context, key string) error
}
