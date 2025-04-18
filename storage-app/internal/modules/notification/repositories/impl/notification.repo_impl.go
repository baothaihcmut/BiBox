package impl

import "go.mongodb.org/mongo-driver/mongo"

type NotificationRepoImpl struct {
	collection *mongo.Collection
}

func NewNotificationRepo(collection *mongo.Collection) *NotificationRepoImpl {
	return &NotificationRepoImpl{
		collection: collection,
	}
}
