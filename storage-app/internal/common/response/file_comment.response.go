package response

import (
	"time"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_comment/models"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileCommentAnswerOutput struct {
	UserId    primitive.ObjectID   `json:"user_id"`
	Content   string               `json:"content"`
	Mentions  []primitive.ObjectID `json:"mentions"`
	CreatedAt time.Time            `json:"created_at"`
}

type FileCommentOutput struct {
	Id        primitive.ObjectID        `json:"id"`
	FileId    primitive.ObjectID        `json:"file_id"`
	UserId    primitive.ObjectID        `json:"user_id"`
	Mentions  []primitive.ObjectID      `json:"mentions"`
	Answers   []FileCommentAnswerOutput `json:"answers"`
	CreatedAt time.Time                 `json:"created_at"`
}

func MapToFileCommentOutput(comment models.FileComment) *FileCommentOutput {
	return &FileCommentOutput{
		Id:        comment.Id,
		FileId:    comment.FileID,
		UserId:    comment.UserID,
		Mentions:  comment.Mentions,
		CreatedAt: comment.CreatedAt,
		Answers: lo.Map(comment.Answers, func(item models.FileCommentAnswer, _ int) FileCommentAnswerOutput {
			return FileCommentAnswerOutput{
				UserId:    item.UserId,
				Mentions:  item.Mentions,
				Content:   item.Content,
				CreatedAt: item.CreatedAt,
			}
		}),
	}

}
