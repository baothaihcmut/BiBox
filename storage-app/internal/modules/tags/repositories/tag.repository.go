package repositories

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/tags/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TagRepository interface {
	CreateTag(context.Context, *models.Tag) error
	BulkCreateTag(context.Context, []*models.Tag) error
	FindTagById(context.Context, primitive.ObjectID) (*models.Tag, error)
	FindAllTagInList(context.Context, []primitive.ObjectID) ([]*models.Tag, error)
	FindAllTagsAndCount(ctx context.Context, query string, limit, offset int) ([]*models.Tag, int, error)
}
