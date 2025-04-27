package presenters

import "github.com/baothaihcmut/Bibox/storage-app/internal/common/response"

type UploadedFileInput struct {
	Id              string `uri:"id" validate:"required"`
	UploadLockValue string `json:"upload_lock_value" validate:"required"`
}

type UploadedFileOutput struct {
	*response.FileOutput
}
