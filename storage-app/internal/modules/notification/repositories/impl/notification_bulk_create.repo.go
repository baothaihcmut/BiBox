package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/notification/models"
	"github.com/samber/lo"
)

func (n *NotificationRepoImpl) BulkCreateNotifications(ctx context.Context, notifications []*models.Notification) error {
	_, err := n.collection.InsertMany(ctx, lo.Map(notifications, func(item *models.Notification, _ int) interface{} {
		return item
	}))
	if err != nil {
		return err
	}
	return nil
}
