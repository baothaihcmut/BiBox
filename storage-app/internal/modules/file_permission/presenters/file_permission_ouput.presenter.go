package presenters

import (
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FilePermissionOuput struct {
	FileID           primitive.ObjectID       `json:"file_id"`
	UserID           primitive.ObjectID       `json:"user_id"`
	PermissionType   enums.FilePermissionType `json:"permission_type"`
	CanShare         bool                     `json:"can_share"`
	AccessSecureFile bool                     `json:"access_secure_file"`
}

func MapToOuput(f models.FilePermission) *FilePermissionOuput {
	return &FilePermissionOuput{
		FileID:         f.FileID,
		UserID:         f.UserID,
		PermissionType: f.PermissionType,
		CanShare:       f.CanShare,
	}
}
