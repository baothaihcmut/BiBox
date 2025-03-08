package models

import (
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileWithPermission struct {
	File        `bson:"inline"`
	Permissions []struct {
		UserID         primitive.ObjectID       `bson:"user_id"`
		PermissionType enums.FilePermissionType `bson:"permission_type"`
		UserImage      string                   `bson:"user_image"`
		FirstName      string                   `bson:"user_first_name"`
		LastName       string                   `bson:"user_last_name"`
		Email          string                   `bson:"user_email"`
	} `bson:"permissions"`
	PermissionType enums.FilePermissionType `bson:"permission_type"`
}
