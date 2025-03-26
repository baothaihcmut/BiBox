package presenters

import (
	"time"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserPermission struct {
	UserId          primitive.ObjectID       `json:"user_id" validate:"required"`
	PermissionsType enums.FilePermissionType `json:"permission_type" validate:"required,gte=1,lte=3"`
	ExpireAt        *time.Time               `json:"expire_at"`
}

type AddFilePermissionInput struct {
	FileId          string           `uri:"id"`
	UserPermissions []UserPermission `json:"permissions" validate:"required"`
}

type AddFilePermissionOutput struct {
	Permissions []*response.FilePermissionOuput `json:"permissions"`
}
