package presenters

import (
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateFileInput struct {
	Name           string               `json:"name" validate:"required"`
	IsFolder       bool                 `json:"is_folder"` // Use *bool to allow nil check
	ParentFolderID *primitive.ObjectID  `json:"parent_folder_id,omitempty"`
	Description    string               `json:"description"`
	TagIDs         []primitive.ObjectID `json:"tags,omitempty"`
	StorageDetail  *struct {
		Size     int    `json:"size" validate:"required"`      // Required field
		MimeType string `json:"mime_type" validate:"required"` // Required field
	} `json:"storage_detail,omitempty"`
}

type CreateFileOutput struct {
	*response.FileOutput

	PutObjectUrl string `json:"put_object_url"`
	UrlExpiry    int    `json:"url_expiry"`
}
