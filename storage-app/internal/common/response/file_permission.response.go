package response

import (
	"time"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FilePermissionOuput struct {
	FileID             primitive.ObjectID       `json:"file_id"`
	UserID             primitive.ObjectID       `json:"user_id"`
	FilePermissionType enums.FilePermissionType `json:"permission_type"`
	CanShare           bool                     `json:"can_share"`
	ExpireAt           *time.Time               `json:"expire_at"`
}

func MapToFilePermissionOutput(f *models.FilePermission) *FilePermissionOuput {
	return &FilePermissionOuput{
		FileID:             f.FileID,
		UserID:             f.UserID,
		FilePermissionType: f.FilePermissionType,
		CanShare:           f.CanShare,
		ExpireAt:           f.ExpireAt,
	}
}
