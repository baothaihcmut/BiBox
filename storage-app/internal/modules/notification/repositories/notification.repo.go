package repositories

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/notification/models"
)

type NotificationRepo interface {
	CreateNotification(ctx context.Context, notification *models.Notification) error
	BulkCreateNotifications(ctx context.Context, notifications []*models.Notification) error
}
