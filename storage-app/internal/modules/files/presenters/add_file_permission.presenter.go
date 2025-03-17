package presenters

import (
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/presenters"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdditionFilePermission struct {
	FileId         primitive.ObjectID       `json:"file_id" binding:"required"`
	PermissionType enums.FilePermissionType `json:"permission_type" binding:"required"`
}

type FilePermission struct {
	UserId             primitive.ObjectID       `json:"user_id" binding:"required"`
	PermissionType     *enums.PermissionType    `json:"permission_type"`
	AdditionPermission []AdditionFilePermission `json:"addition_permissions"`
}

type AddFilePermissionInput struct {
	FileId      string           `uri:"id"`
	Permissions []FilePermission `json:"permissions"`
}

type AddFilePermissionOutput struct {
	Permissions []*presenters.FilePermissionOuput `json:"permissions"`
}
