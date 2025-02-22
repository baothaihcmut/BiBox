package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileComment struct {
	FileID    primitive.ObjectID `bson:"file_id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id,omitempty"`
	Content   string             `bson:"content,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty"`
	Answers   []string           `bson:"answers,omitempty"` // Nested comments
}
