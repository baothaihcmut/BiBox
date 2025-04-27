package presenters

import "github.com/baothaihcmut/Bibox/storage-app/internal/common/response"

type UpdateFileContentInput struct {
	Id            string `uri:"id" validate:"required"`
	StorageDetail *struct {
		Size     int    `json:"size" validate:"gt=0"`
		MimeType string `json:"mime_type"`
	} `json:"storage_detail"`
}
type UpdateFileContentOutput struct {
	*response.FileOutput
	PutObjectUrl    string `json:"put_object_url"`
	UrlExpiry       int    `json:"url_expiry"`
	UploadLockValue string `json:"upload_lock_value"`
}
