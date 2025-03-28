package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/models"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/repositories"
	"go.mongodb.org/mongo-driver/bson"
)

func (pr *FilePermissionRepositoryImpl) FindPermissionByListId(ctx context.Context, ids []repositories.FilePermissionId) ([]*models.FilePermission, error) {
	orFilterIds := bson.A{}
	for _, id := range ids {
		orFilterIds = append(orFilterIds, bson.M{
			"file_id": id.FileId,
			"user_id": id.UserId,
		})
	}
	cursor, err := pr.collection.Find(ctx, bson.M{
		"$or": orFilterIds,
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
