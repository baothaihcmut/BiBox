package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_comment/models"
)

func (f *MongoFileCommentRepo) CreateComment(ctx context.Context, comment *models.FileComment) error {
	_, err := f.collection.InsertOne(ctx, comment)
	if err != nil {
		return err
	}
	return nil
}
