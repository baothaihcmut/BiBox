package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/models"
	"github.com/samber/lo"
)

func (pr *FilePermissionRepositoryImpl) BulkCreatePermission(ctx context.Context, filePermissions []*models.FilePermission) error {
	_, err := pr.collection.InsertMany(ctx, lo.Map(filePermissions, func(item *models.FilePermission, _ int) interface{} {
		return item
	}))
	if err != nil {
		return err
	}
	return nil

}
