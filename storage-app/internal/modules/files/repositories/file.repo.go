package repositories

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FindFileWithPermissionArg struct {
	IsFolder       *bool
	SortBy         string
	IsAsc          bool
	Offset         int
	Limit          int
	UserId         primitive.ObjectID //user context
	ParentFolderId *primitive.ObjectID
	FileType       *enums.MimeType
	OwnerId        *primitive.ObjectID
	IsDeleted      *bool
	TagId          *primitive.ObjectID
}

type FileRepository interface {
	CreateFile(context.Context, *models.File) error
	BulkCreateFiles(context.Context, []*models.File) error
	BulkUpdateFile(context.Context, []*models.File) error
	UpdateFile(context.Context, *models.File) error
	BulkDeleteFile(context.Context, []*models.File) error

	FindFileById(context.Context, primitive.ObjectID) (*models.File, error)
	FindSubFileRecursive(context.Context, primitive.ObjectID) ([]*models.File, error)
	FindFileByParentFolderId(context.Context, primitive.ObjectID) ([]*models.File, error)
	FindAllParentFolder(context.Context, primitive.ObjectID) ([]*models.File, error)

	FindFileWithPermssionAndCount(context.Context, FindFileWithPermissionArg) ([]*models.FileWithPermission, int64, error)
}
