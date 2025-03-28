package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (pr *FilePermissionRepositoryImpl) DeletePermission(ctx context.Context, filePermission *models.FilePermission) error {
	_, err := pr.collection.DeleteOne(ctx, bson.M{
		"file_id": filePermission.FileID,
		"user_id": filePermission.UserID,
	})
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	return nil
}
