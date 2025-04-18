package impl

import (
	"context"
	"fmt"
	"sync"

	"github.com/baothaihcmut/BiBox/libs/pkg/events/notifications"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/notification/handlers"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/notification/models"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/notification/services"
	"github.com/google/uuid"
)

func (n *NotificationServiceImpl) SendNotificationFileUploaded(ctx context.Context, args []services.SendNotificationFileUploadedArg) error {

	newNotifications := make([]*models.Notification, 0, len(args))
	for _, arg := range args {
		newNotifications = append(newNotifications, models.NewNotification(
			arg.FileOwnerId,
			enums.FileUploaded,
			"User uploaded file",
			"New user has uploaded file in your folder",
			fmt.Sprintf("localhost:3000/files/%s", arg.FileId.Hex()),
			arg.UserUploadId,
		))
	}
	if err := n.repo.BulkCreateNotifications(ctx, newNotifications); err != nil {
		return err
	}
	//publish to queue
	wg := &sync.WaitGroup{}
	for _, notification := range newNotifications {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _, err := n.queueService.PublishMessage(
				ctx,
				handlers.NotificationCreatedTopic,
				&notifications.NotificationCreatedEvent{
					Id:         notification.Id.Hex(),
					UserId:     notification.UserId.Hex(),
					FromUserId: notification.FromUserId.Hex(),
					Type:       int(notification.Type),
					ActionUrl:  notification.ActionUrl,
					Title:      notification.Title,
					Message:    notification.Message,
					CreatedAt:  notification.CreatedAt,
					Seen:       notification.Seen,
				}, map[string]string{
					"eventType":   "FileUploaded",
					"eventSource": "storage-app",
					"enventId":    uuid.NewString(),
				},
			)
			if err != nil {
				n.logger.Errorf(ctx, nil, "Error publish to queue:", err)
			}
		}()
	}
	wg.Wait()
	return nil
}
