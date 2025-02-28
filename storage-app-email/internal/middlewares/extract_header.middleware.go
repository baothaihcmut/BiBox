package middlewares

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/baothaihcmut/BiBox/storage-app-email/internal/constant"
	"github.com/baothaihcmut/BiBox/storage-app-email/internal/router"
)

func ExtractHeaderMiddleware(handler router.HandleFunc) router.HandleFunc {
	return func(ctx context.Context, cm *sarama.ConsumerMessage) error {
		//extract header
		headers := make(map[string]string)
		for _, header := range cm.Headers {
			headers[string(header.Key)] = string(header.Value)
		}
		ctx = context.WithValue(ctx, constant.HeaderContext, headers)
		return handler(ctx, cm)
	}
}
