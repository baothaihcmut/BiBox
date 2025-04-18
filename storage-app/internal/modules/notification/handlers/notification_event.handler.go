package handlers

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/baothaihcmut/BiBox/libs/pkg/events/notifications"
	"github.com/baothaihcmut/BiBox/libs/pkg/middlewares"
	"github.com/baothaihcmut/BiBox/libs/pkg/router"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/constant"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/notification/services"
)

const NotificationCreatedTopic = "notifications.created"

type NotificationEventHandler interface {
	Init(r router.MessageRouter)
}

type KafkaNotificationEventHandler struct {
	msgChs                        map[string]chan *sarama.ConsumerMessage
	notificationSSEManagerService services.NotificationSSEManagerService
}

func (k *KafkaNotificationEventHandler) handleNotificationCreatedEvent(ctx context.Context, msg *sarama.ConsumerMessage) error {
	e, _ := ctx.Value(constant.PayloadContext).(*notifications.NotificationCreatedEvent)
	err := k.notificationSSEManagerService.SendNotificationCreatedEvent(ctx, e)
	if err != nil {
		return err
	}
	return nil
}

func (k *KafkaNotificationEventHandler) Init(r router.MessageRouter) {
	r.Register(NotificationCreatedTopic, k.handleNotificationCreatedEvent, middlewares.ExtractEventMiddleware[notifications.NotificationCreatedEvent]())
}
func NewNotificationEventHandler(notificationService services.NotificationSSEManagerService) *KafkaNotificationEventHandler {
	return &KafkaNotificationEventHandler{
		msgChs:                        make(map[string]chan *sarama.ConsumerMessage),
		notificationSSEManagerService: notificationService,
	}
}
