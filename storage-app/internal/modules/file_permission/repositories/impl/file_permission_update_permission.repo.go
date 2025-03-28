package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (pr *FilePermissionRepositoryImpl) UpdatePermission(ctx context.Context, filePermission *models.FilePermission) error {
	_, err := pr.collection.UpdateOne(ctx, bson.M{
		"file_id": filePermission.FileID,
		"user_id": filePermission.UserID,
	},
		bson.M{
			"$set": bson.M{
				"permission_type": filePermission.FilePermissionType,
				"expire_at":       filePermission.ExpireAt,
			},
		},
	)
	if err != nil {
		return err
	}
	return nil
}
