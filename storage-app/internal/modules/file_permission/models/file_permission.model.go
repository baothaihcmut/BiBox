package models

import (
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FilePermission struct {
	FileID           primitive.ObjectID       `bson:"file_id"`
	UserID           primitive.ObjectID       `bson:"user_id"`
	PermissionType   enums.FilePermissionType `bson:"permission_type"`
	CanShare         bool                     `bson:"can_share"`
	AccessSecureFile bool                     `bson:"access_secure_file"`
}

func NewFilePermission(fileID, userID primitive.ObjectID, permissionType enums.FilePermissionType, canShare, accessSecure bool) *FilePermission {
	return &FilePermission{
		FileID:           fileID,
		UserID:           userID,
		PermissionType:   permissionType,
		CanShare:         canShare,
		AccessSecureFile: accessSecure,
	}
}
