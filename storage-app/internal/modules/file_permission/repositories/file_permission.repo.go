package repositories

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FilterPermssionOption int

type FilePermissionId struct {
	FileId primitive.ObjectID
	UserId primitive.ObjectID
}

type FilePermissionRepository interface {
	CreateFilePermission(ctx context.Context, filePermission *models.FilePermission) error
	BulkCreatePermission(ctx context.Context, filePermissions []*models.FilePermission) error
	BulkUpdatePermission(ctx context.Context, permissions []*models.FilePermission) error
	UpdatePermission(ctx context.Context, filePermission *models.FilePermission) error
	DeletePermission(ctx context.Context, filePermission *models.FilePermission) error
	BulkDeletePermission(ctx context.Context, filePermissions []*models.FilePermission) error
	DeletePermissionByListFileId(ctx context.Context, fileIds []primitive.ObjectID) error

	FindFilePermissionById(ctx context.Context, id FilePermissionId) (*models.FilePermission, error)
	FindPermissionByListId(ctx context.Context, ids []FilePermissionId) ([]*models.FilePermission, error)
	FindPermssionByFileId(ctx context.Context, fileId primitive.ObjectID) ([]*models.FilePermission, error)

	FindFilePermissionWithUser(ctx context.Context, fileId primitive.ObjectID) ([]*models.FilePermissionWithUser, error)
}
