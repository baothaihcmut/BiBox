package presenters

import (
	"time"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/response"
)

type UpdateFilePermissionInput struct {
	FileId         string                   `uri:"id" validate:"required"`
	UserId         string                   `uri:"userId" validate:"required"`
	PermissionType enums.FilePermissionType `json:"permission_type" validate:"required,gte=1,lte=3"`
	ExpireAt       *time.Time               `json:"expire_at"`
}

type UpdateFilePermissionOuput struct {
	Permissions []*response.FilePermissionOuput `json:"permissions"`
}
