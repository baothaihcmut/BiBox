package models

import (
	"time"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FilePermission struct {
	FileID         primitive.ObjectID       `bson:"file_id"`
	UserID         primitive.ObjectID       `bson:"user_id"`
	PermissionType enums.FilePermissionType `bson:"permission_type"`
	CanShare       bool                     `bson:"can_share"`
	ExpireAt       *time.Time               `bson:"expire_at"`
}

func NewFilePermission(fileID, userID primitive.ObjectID, permissionType enums.FilePermissionType, canShare bool, exprireAt *time.Time) *FilePermission {
	return &FilePermission{
		FileID:         fileID,
		UserID:         userID,
		PermissionType: permissionType,
		CanShare:       canShare,
		ExpireAt:       exprireAt,
	}
}
