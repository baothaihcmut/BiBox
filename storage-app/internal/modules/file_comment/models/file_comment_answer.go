package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileCommentAnswer struct {
	UserId    primitive.ObjectID   `bson:"user_id"`
	Content   string               `bson:"content"`
	Mentions  []primitive.ObjectID `bson:"mentions"`
	CreatedAt time.Time            `bson:"created_at"`
}

func NewFileCommentAnswer(
	userId primitive.ObjectID,
	content string,
	mentions []primitive.ObjectID,
) *FileCommentAnswer {
	return &FileCommentAnswer{
		UserId:    userId,
		Content:   content,
		Mentions:  mentions,
		CreatedAt: time.Now(),
	}
}
