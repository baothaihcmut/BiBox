package lock

import (
	"context"
	"time"
)

type DistributedLockService interface {
	AcquireLock(ctx context.Context, key string, value string, maxRetry int, interval time.Duration, ttl time.Duration) (bool, error)
	ReleaseLock(ctx context.Context, key string, value string) error
}
