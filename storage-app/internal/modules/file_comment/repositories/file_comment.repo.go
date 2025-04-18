package repositories

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_comment/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CommentRepo
type FileCommentRepository interface {
	CreateComment(ctx context.Context, comment *models.FileComment) error
	FindCommentById(ctx context.Context, id primitive.ObjectID) (*models.FileComment, error)
	FindCommentByFileId(ctx context.Context, fileId primitive.ObjectID) ([]*models.FileComment, error)
}
