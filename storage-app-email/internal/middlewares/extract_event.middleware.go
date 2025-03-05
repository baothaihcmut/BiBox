package middlewares

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/baothaihcmut/BiBox/storage-app-email/internal/constant"
	"github.com/baothaihcmut/BiBox/storage-app-email/internal/router"
)

func ExtractEventMiddleware[T any]() router.MiddlewareFunc {
	return func(handler router.HandleFunc) router.HandleFunc {
		return func(ctx context.Context, cm *sarama.ConsumerMessage) error {
			//extract event
			var e T
			if err := json.Unmarshal(cm.Value, &e); err != nil {
				return err
			}
			ctx = context.WithValue(ctx, constant.PayloadContext, &e)
			return handler(ctx, cm)
		}
	}
}
