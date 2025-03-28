package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/models"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (pr *FilePermissionRepositoryImpl) FindFilePermissionById(ctx context.Context, id repositories.FilePermissionId) (*models.FilePermission, error) {
	// build filter
	var file models.FilePermission
	err := pr.collection.FindOne(ctx, bson.M{
		"file_id": id.FileId,
		"user_id": id.UserId,
	}).Decode(&file)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // File permission not found
		}
		return nil, err
	}
	return &file, nil
}
