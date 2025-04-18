package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/notification/models"
)

func (n *NotificationRepoImpl) CreateNotification(ctx context.Context, notification *models.Notification) error {
	_, err := n.collection.InsertOne(ctx, notification)
	if err != nil {
		return err
	}
	return nil

}
