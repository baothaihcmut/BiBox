package initialize

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/config"
	"github.com/redis/go-redis/v9"
)

func InitializeRedis(cfg *config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Endpoint,
		Username: cfg.Username,
		Password: cfg.Password,
		DB:       cfg.Database,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}
