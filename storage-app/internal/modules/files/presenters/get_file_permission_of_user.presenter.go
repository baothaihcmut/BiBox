package presenters

import "github.com/baothaihcmut/Bibox/storage-app/internal/common/response"

type GetFilePermissionOfUserInput struct {
	FileId string `uri:"id" validate:"required"`
}

type GetFilePermissionOfUserOutput struct {
	response.FilePermissionOuput
}
