package presenters

import (
	"time"
)

type CreateFileInput struct {
	IsFolder       bool     `json:"is_folder"`
	ParentFolderID *string  `json:"parent_folder_id,omitempty"`
	HasPassword    bool     `json:"has_password"`
	Password       *string  `json:"password,omitempty"`
	Description    string   `json:"description"`
	IsSecure       bool     `json:"is_secure"`
	TagIDs         []string `json:"tags,omitempty"`
	StorageDetail  *struct {
		Size int    `json:"size"`
		Type string `json:"file_type"`
	} `json:"storage_detail,omitempty"`
}

type CreateFileOutput struct {
	ID             string              `json:"id"`
	OwnerID        string              `json:"owner_id"`
	IsFolder       bool                `json:"is_folder"`
	ParentFolderID *string             `json:"parent_folder_id"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
	OpenedAt       *time.Time          `json:"opened_at"`
	HasPassword    bool                `json:"has_password"`
	Description    string              `json:"description"`
	IsDeleted      bool                `json:"is_deleted"`
	DeletedAt      *time.Time          `json:"deleted_at"`
	TotalSize      int                 `json:"total_size"`
	IsSecure       bool                `json:"is_secure"`
	TagIDs         []string            `json:"tags"`
	StorageDetails *StorageDetailOuput `json:"storage_detail"`
}
type StorageDetailOuput struct {
	Size         int    `json:"file_size"`
	Type         string `json:"file_type"`
	PutObjectUrl string `json:"put_object_url"`
	UrlExpiry    int    `json:"url_expiry"`
}
