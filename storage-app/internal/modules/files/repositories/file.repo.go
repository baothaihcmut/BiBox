package repositories

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FindFileOfUserArg struct {
	IsFolder       *bool
	ParentFolderId *primitive.ObjectID
	SortBy         string
	IsAsc          bool
	Offset         int
	Limit          int
	PermssionLimit int
	FileType       *enums.MimeType
	OwnerId        *primitive.ObjectID
	UserId         primitive.ObjectID
}

type FileRepository interface {
	CreateFile(context.Context, *models.File) error
	BulkCreateFiles(context.Context, []*models.File) error
	BulkUpdateTotalSize(context.Context, []*models.File) error
	UpdateFile(context.Context, *models.File) error
	FindFileById(context.Context, primitive.ObjectID, bool) (*models.File, error)
	FindFileWithPermssionAndCount(context.Context, FindFileOfUserArg) ([]*models.FileWithPermission, int64, error)
	FindSubFileRecursive(context.Context, primitive.ObjectID) ([]*models.File, error)
	FindFileByParentFolderId(context.Context, primitive.ObjectID) ([]*models.File, error)
	FindAllParentFolder(context.Context, primitive.ObjectID) ([]*models.File, error)
	FindAllFileByTagAndCount(
		ctx context.Context,
		tagId, userId primitive.ObjectID,
		limit, offset int,
		sortBy string,
		isAsc bool,
	) ([]*models.FileWithPermission, int, error)
}
