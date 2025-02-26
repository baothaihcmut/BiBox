package middlewares

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/baothaihcmut/BiBox/storage-app-email/internal/constant"
	"github.com/baothaihcmut/BiBox/storage-app-email/internal/router"
)

func ExtractEventMiddleware(handler router.HandleFunc) router.HandleFunc {
	return func(ctx context.Context, cm *sarama.ConsumerMessage) error {
		var e interface{}
		if err := json.Unmarshal(cm.Value, &e); err != nil {
			return err
		}
		ctx = context.WithValue(ctx, constant.PayloadContext, &e)
		return handler(ctx, cm)
	}
}
