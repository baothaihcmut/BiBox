package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (pr *FilePermissionRepositoryImpl) BulkUpdatePermission(ctx context.Context, permissions []*models.FilePermission) error {
	bulkOperation := make([]mongo.WriteModel, 0, len(permissions))
	for _, permission := range permissions {
		filter := bson.M{
			"file_id": permission.FileID,
			"user_id": permission.UserID,
		}
		update := bson.M{
			"$set": bson.M{
				"permission_type": permission.FilePermissionType,
				"expire_at":       permission.ExpireAt,
				"can_share":       permission.CanShare,
			},
		}
		bulkOperation = append(bulkOperation, mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update))
	}
	if len(bulkOperation) > 0 {
		_, err := pr.collection.BulkWrite(ctx, bulkOperation)
		if err != nil {
			return err
		}
	}
	return nil
}
