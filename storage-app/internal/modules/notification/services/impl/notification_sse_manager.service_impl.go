package impl

import (
	"context"

	"github.com/baothaihcmut/BiBox/libs/pkg/events/notifications"
	"github.com/baothaihcmut/BiBox/libs/pkg/logger"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/cache"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/response"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/sse"
)

const notificationSessionKey = "session:notification"

type NotificationSSEManagerService struct {
	*sse.SSEManagerService[response.NotificationOutput]
}

func NewNotificationSSEManagerService(
	cacheService cache.CacheService,
	logger logger.Logger,
) *NotificationSSEManagerService {
	return &NotificationSSEManagerService{
		sse.NewNotificationSSEManagerService[response.NotificationOutput](
			cacheService, logger,
		),
	}
}

func (n *NotificationSSEManagerService) Connect(ctx context.Context, sessionId string) (<-chan *response.NotificationOutput, string, error) {
	return n.SSEManagerService.Connect(ctx, notificationSessionKey, sessionId)
}

func (n *NotificationSSEManagerService) SendNotificationCreatedEvent(ctx context.Context, e *notifications.NotificationCreatedEvent) error {
	msg := &response.NotificationOutput{
		Id:         e.Id,
		UserId:     e.UserId,
		Type:       enums.NoficationType(e.Type),
		Title:      e.Title,
		Message:    e.Message,
		ActionUrl:  e.ActionUrl,
		Seen:       e.Seen,
		FromUserId: e.FromUserId,
		CreatedAt:  e.CreatedAt,
	}
	n.SSEManagerService.SendMessage(ctx, e.UserId, msg)
	return nil
}

func (n *NotificationSSEManagerService) ClearClosedSession(ctx context.Context) {
	n.SSEManagerService.ClearClosedSession(ctx, notificationSessionKey)
}
func (n *NotificationSSEManagerService) Disconnect(ctx context.Context, notificationId string, sessionId string) error {
	return n.SSEManagerService.Disconnect(ctx, notificationId, sessionId)
}
