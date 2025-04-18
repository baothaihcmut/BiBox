package services

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SendNotificationFileUploadedArg struct {
	FileOwnerId  primitive.ObjectID
	UserUploadId primitive.ObjectID
	FileId       primitive.ObjectID
}

type NotificationService interface {
	SendNotificationFileUploaded(ctx context.Context, args []SendNotificationFileUploadedArg) error
}
