package models

import "time"

type FileComment struct {
	FileID    string    `bson:"file_id,omitempty"`
	UserID    string    `bson:"user_id,omitempty"`
	Content   string    `bson:"content,omitempty"`
	CreatedAt time.Time `bson:"created_at,omitempty"`
	Answers   []string  `bson:"answers,omitempty"` // Nested comments
}
