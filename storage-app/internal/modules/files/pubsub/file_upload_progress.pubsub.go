package handlers

import (
	"context"
	"encoding/json"

	"github.com/baothaihcmut/BiBox/libs/pkg/logger"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/services"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/services/impl"
	"github.com/redis/go-redis/v9"
)

type FileUploadProgressPubSub interface {
	Run(context.Context, string)
}

type RedisFileUploadProgressPubSub struct {
	client                   *redis.Client
	uploadProgressSSEManager services.FileUploadProgressSSEManagerService
	logger                   logger.Logger
}

func (r *RedisFileUploadProgressPubSub) Run(ctx context.Context, channel string) {
	pubsub := r.client.Subscribe(ctx, channel)
	ch := pubsub.Channel()
	for msg := range ch {
		var payload impl.FileUploadProgressPayload
		if err := json.Unmarshal([]byte(msg.Payload), &payload); err != nil {
			r.logger.Errorf(ctx, nil, "Error encode payload:", err)
		}
		msg := &services.FileUploadProgress{
			UploadSpeed: payload.UploadSpeed,
			TotalSize:   payload.TotalSize,
			Percent:     payload.Percent,
		}
		if err := r.uploadProgressSSEManager.SendUploadProgressMessage(ctx, payload.FileId, msg); err != nil {
			r.logger.Errorf(ctx, nil, "Error send message to sse")
		}
	}
}
func NewRedisFileUploadProgressPubSub(
	client *redis.Client,
	sseManager services.FileUploadProgressSSEManagerService,
	logger logger.Logger,
) *RedisFileUploadProgressPubSub {
	return &RedisFileUploadProgressPubSub{
		client:                   client,
		uploadProgressSSEManager: sseManager,
		logger:                   logger,
	}
}
