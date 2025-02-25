package presenters

import (
	"time"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/models"
)

type StorageDetailOuput struct {
	Size     int            `json:"file_size"`
	MimeType enums.MimeType `json:"mime_type"`
}
type FileOutput struct {
	ID             string              `json:"id"`
	Name           string              `json:"name"`
	OwnerID        string              `json:"owner_id"`
	IsFolder       bool                `json:"is_folder"`
	ParentFolderID *string             `json:"parent_folder_id"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
	OpenedAt       *time.Time          `json:"opened_at"`
	HasPassword    bool                `json:"has_password"`
	Description    string              `json:"description"`
	IsSecure       bool                `json:"is_secure"`
	TagIDs         []string            `json:"tags"`
	StorageDetails *StorageDetailOuput `json:"storage_detail"`
}

func MapFileToFileOutput(file *models.File) *FileOutput {
	// Convert primitive.ObjectID to string
	var parentFolderID *string
	if file.ParentFolderID != nil {
		id := file.ParentFolderID.Hex()
		parentFolderID = &id
	}

	tagIDs := make([]string, len(file.TagIDs))
	for i, tagID := range file.TagIDs {
		tagIDs[i] = tagID.Hex()
	}

	var storageDetailOutput *StorageDetailOuput
	if file.StorageDetail != nil {
		storageDetailOutput = &StorageDetailOuput{
			Size:     file.StorageDetail.Size,
			MimeType: file.StorageDetail.MimeType,
		}
	}

	return &FileOutput{
		ID:             file.ID.Hex(),
		Name:           file.Name,
		OwnerID:        file.OwnerID.Hex(),
		IsFolder:       file.IsFolder,
		ParentFolderID: parentFolderID,
		CreatedAt:      file.CreatedAt,
		UpdatedAt:      file.UpdatedAt,
		OpenedAt:       file.OpenedAt,
		HasPassword:    file.HasPassword,
		Description:    file.Description,
		IsSecure:       file.IsSecure,
		TagIDs:         tagIDs,
		StorageDetails: storageDetailOutput,
	}
}
