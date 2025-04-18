package models

import (
	"time"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Notification struct {
	Id         primitive.ObjectID   `bson:"_id"`
	UserId     primitive.ObjectID   `bson:"user_id"`
	Type       enums.NoficationType `bson:"type"`
	Title      string               `bson:"title"`
	Message    string               `bson:"message"`
	ActionUrl  string               `bson:"action_url"`
	Seen       bool                 `bson:"seen"`
	FromUserId primitive.ObjectID   `bson:"from_user_id"`
	CreatedAt  time.Time            `bson:"created_at"`
}

func NewNotification(
	userId primitive.ObjectID,
	notiType enums.NoficationType,
	title string,
	message string,
	actionUrl string,
	fromUserId primitive.ObjectID,
) *Notification {
	return &Notification{
		Id:         primitive.NewObjectID(),
		UserId:     userId,
		Type:       notiType,
		Title:      title,
		Message:    message,
		ActionUrl:  actionUrl,
		Seen:       false,
		FromUserId: fromUserId,
		CreatedAt:  time.Now(),
	}
}
