package presenters

import (
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateFileInput struct {
	Name           string               `json:"name" validate:"required"`
	IsFolder       bool                 `json:"is_folder"`
	ParentFolderID *primitive.ObjectID  `json:"parent_folder_id"`
	Description    string               `json:"description"`
	TagIDs         []primitive.ObjectID `json:"tags" validate:"required"`
	StorageDetail  *struct {
		Size     int    `json:"size" validate:"gt=0"`
		MimeType string `json:"mime_type"`
	} `json:"storage_detail"`
}

type CreateFileOutput struct {
	*response.FileOutput
	PutObjectUrl    string `json:"put_object_url"`
	UrlExpiry       int    `json:"url_expiry"`
	UploadLockValue string `json:"upload_lock_value"`
}
