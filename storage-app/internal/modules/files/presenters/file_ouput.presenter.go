package presenters

import (
	"time"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileWithPermissionOutput struct {
	*FileOutput
	Permissions    []*PermissionOfFileOuput `json:"permissions"`
	PermissionType enums.FilePermissionType `json:"permission_type"`
}

type StorageDetailOuput struct {
	Size     int            `json:"file_size"`
	MimeType enums.MimeType `json:"mime_type"`
}
type FileOutput struct {
	ID             primitive.ObjectID   `json:"id"`
	Name           string               `json:"name"`
	OwnerID        primitive.ObjectID   `json:"owner_id"`
	IsFolder       bool                 `json:"is_folder"`
	ParentFolderID *primitive.ObjectID  `json:"parent_folder_id"`
	CreatedAt      time.Time            `json:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at"`
	OpenedAt       *time.Time           `json:"opened_at"`
	HasPassword    bool                 `json:"has_password"`
	Description    string               `json:"description"`
	IsSecure       bool                 `json:"is_secure"`
	TagIDs         []primitive.ObjectID `json:"tags"`
	StorageDetails *StorageDetailOuput  `json:"storage_detail"`
}

type FileOwnerOuput struct {
	Id        primitive.ObjectID `json:"id"`
	Email     string             `json:"email"`
	FirstName string             `json:"first_name"`
	LastName  string             `json:"last_name"`
	Image     string             `json:"image"`
}

type PermissionOfFileOuput struct {
	UserID         primitive.ObjectID       `json:"user_id"`
	PermissionType enums.FilePermissionType `json:"permission_type"`
	UserImage      string                   `json:"user_image"`
}

func MapFileToFileOutput(file *models.File) *FileOutput {

	var storageDetailOutput *StorageDetailOuput
	if file.StorageDetail != nil {
		storageDetailOutput = &StorageDetailOuput{
			Size:     file.StorageDetail.Size,
			MimeType: file.StorageDetail.MimeType,
		}
	}

	return &FileOutput{
		ID:             file.ID,
		Name:           file.Name,
		OwnerID:        file.OwnerID,
		IsFolder:       file.IsFolder,
		ParentFolderID: file.ParentFolderID,
		CreatedAt:      file.CreatedAt,
		UpdatedAt:      file.UpdatedAt,
		OpenedAt:       file.OpenedAt,
		HasPassword:    file.HasPassword,
		Description:    file.Description,
		IsSecure:       file.IsSecure,
		TagIDs:         file.TagIDs,
		StorageDetails: storageDetailOutput,
	}
}
