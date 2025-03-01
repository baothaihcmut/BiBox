package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type FilePermission struct {
	FileID         primitive.ObjectID `bson:"file_id"`
	UserID         primitive.ObjectID `bson:"user_id"`
	PermissionType string             `bson:"permission_type"`
	OwnerUserID    primitive.ObjectID `bson:"owner_user_id"`
	CanShare       bool               `bson:"can_share"`
}
