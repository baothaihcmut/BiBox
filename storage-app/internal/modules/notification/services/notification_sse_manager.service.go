package services

import (
	"context"

	"github.com/baothaihcmut/BiBox/libs/pkg/events/notifications"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/response"
)

type NotificationSSEManagerService interface {
	Connect(ctx context.Context, sessionId string) (<-chan *response.NotificationOutput, string, error)
	ClearClosedSession(ctx context.Context)
	SendNotificationCreatedEvent(ctx context.Context, e *notifications.NotificationCreatedEvent) error
	Disconnect(ctx context.Context, notificationId string, sessionId string) error
}
