package middlewares

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/baothaihcmut/BiBox/storage-app-email/internal/constant"
	"github.com/baothaihcmut/BiBox/storage-app-email/internal/router"
)

func LoggingMiddleware(hf router.HandleFunc) router.HandleFunc {
	return func(ctx context.Context, cm *sarama.ConsumerMessage) error {
		headers := ctx.Value(constant.HeaderContext).(map[string]string)
		fmt.Printf(
			"Consume message {event_id: %s,event_source: %s}\n",
			headers["event_id"],
			headers["event_source"],
		)
		err := hf(ctx, cm)
		if err != nil {
			fmt.Printf(
				"Error {event_id: %s,event_source: %s,err:%v}\n",
				headers["event_id"],
				headers["event_source"],
				err,
			)
		} else {
			fmt.Printf(
				"success {event_id: %s,event_source: %s}\n",
				headers["event_id"],
				headers["event_source"],
			)
		}
		return err
	}
}
