package presenters

import "github.com/baothaihcmut/Storage-app/internal/common/enums"

type CreateFileInput struct {
	Name           string   `json:"name" validate:"required"`
	IsFolder       bool     `json:"is_folder"` // Use *bool to allow nil check
	ParentFolderID *string  `json:"parent_folder_id,omitempty"`
	HasPassword    bool     `json:"has_password"` // Use *bool to allow nil check
	Password       *string  `json:"password,omitempty"`
	Description    string   `json:"description"`
	IsSecure       bool     `json:"is_secure"` // Use *bool to allow nil check
	TagIDs         []string `json:"tags,omitempty"`
	StorageDetail  *struct {
		Size     int            `json:"size" validate:"required"`      // Required field
		MimeType enums.MimeType `json:"mime_type" validate:"required"` // Required field
	} `json:"storage_detail,omitempty"`
}

type CreateFileOutput struct {
	*FileOutput
	PutObjectUrl string `json:"put_object_url"`
	UrlExpiry    int    `json:"url_expiry"`
}
