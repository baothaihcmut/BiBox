package presenters

import (
	"time"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserPermission struct {
	UserId          primitive.ObjectID       `json:"user_id" binding:"required"`
	PermissionsType enums.FilePermissionType `json:"permission_type" binding:"required,enum"`
	ExpireAt        *time.Time               `json:"expire_at"`
}

type AddFilePermissionInput struct {
	FileId          string `uri:"id"`
	UserPermissions []UserPermission
}

type AddFilePermissionOutput struct {
	Permissions []*response.FilePermissionOuput `json:"permissions"`
}
