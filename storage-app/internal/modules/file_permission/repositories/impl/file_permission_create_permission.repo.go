package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/models"
)

func (pr *FilePermissionRepositoryImpl) CreateFilePermission(ctx context.Context, filePermission *models.FilePermission) error {
	_, err := pr.collection.InsertOne(ctx, filePermission)
	if err != nil {
		pr.logger.Errorf(ctx, map[string]any{
			"file_id": filePermission.FileID,
			"user_id": filePermission.UserID,
		}, "Error when insert file permission: ", err)
		return err
	}
	return nil
}
