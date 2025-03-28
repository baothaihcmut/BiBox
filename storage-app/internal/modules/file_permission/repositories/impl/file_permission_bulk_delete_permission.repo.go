package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (pr *FilePermissionRepositoryImpl) BulkDeletePermission(ctx context.Context, filePermissions []*models.FilePermission) error {
	deleteIds := bson.A{}
	for _, permission := range filePermissions {
		deleteIds = append(deleteIds, bson.M{
			"file_id": permission.FileID,
			"user_id": permission.UserID,
		})
	}
	_, err := pr.collection.DeleteMany(ctx, bson.M{
		"$or": deleteIds,
	})
	if err != nil {
		return err
	}
	return nil
}
