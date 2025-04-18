package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateFileCommentAnswerArg struct {
	UserId   primitive.ObjectID
	Content  string
	Mentions []primitive.ObjectID
}

type FileComment struct {
	Id        primitive.ObjectID   `bson:"_id"`
	FileID    primitive.ObjectID   `bson:"file_id"`
	UserID    primitive.ObjectID   `bson:"user_id"`
	Content   string               `bson:"content"`
	Mentions  []primitive.ObjectID `bson:"mentions"`
	CreatedAt time.Time            `bson:"created_at"`
	Answers   []FileCommentAnswer  `bson:"answers"`
}

func NewFileComment(
	fileId primitive.ObjectID,
	userId primitive.ObjectID,
	content string,
	mentions []primitive.ObjectID,
) *FileComment {
	return &FileComment{
		Id:        primitive.NewObjectID(),
		FileID:    fileId,
		UserID:    userId,
		Content:   content,
		CreatedAt: time.Now(),
		Answers:   make([]FileCommentAnswer, 0),
		Mentions:  mentions,
	}
}

func (f *FileComment) AddAnswer(arg CreateFileCommentAnswerArg) {
	f.Answers = append(f.Answers, FileCommentAnswer{
		UserId:    arg.UserId,
		Mentions:  arg.Mentions,
		Content:   arg.Content,
		CreatedAt: time.Now(),
	})
}
