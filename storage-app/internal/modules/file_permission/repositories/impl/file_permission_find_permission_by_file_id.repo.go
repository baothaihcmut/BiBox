package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (f *FilePermissionRepositoryImpl) FindPermssionByFileId(ctx context.Context, fileId primitive.ObjectID) ([]*models.FilePermission, error) {
	cursor, err := f.collection.Find(ctx, bson.M{
		"file_id": fileId,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var res []*models.FilePermission
	if err := cursor.All(ctx, &res); err != nil {
		return nil, err
	}
	return res, nil
}
