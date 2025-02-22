package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type FilePermission struct {
	FileID         primitive.ObjectID `bson:"file_id,omitempty"`
	UserID         primitive.ObjectID `bson:"user_id,omitempty"`
	PermissionType string             `bson:"permission_type,omitempty"` 
	AccessSecure   bool               `bson:"access_secure_file"`
}
