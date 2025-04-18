package queue

import "context"

type QueueService interface {
	PublishMessage(ctx context.Context, topic string, value interface{}, headers map[string]string) (int32, int64, error)
}
