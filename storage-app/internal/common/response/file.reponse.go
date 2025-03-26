package response

import (
	"time"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileWithPermissionOutput struct {
	*FileOutput
	Permissions        []*PermissionOfFileOuput `json:"permissions"`
	FilePermissionType enums.FilePermissionType `json:"permission_type"`
}

type StorageDetailOuput struct {
	Size     int            `json:"file_size"`
	MimeType enums.MimeType `json:"mime_type"`
}
type FileOutput struct {
	ID             primitive.ObjectID   `json:"id"`
	Name           string               `json:"name"`
	TotalSize      int                  `json:"total_size"`
	OwnerID        primitive.ObjectID   `json:"owner_id"`
	IsFolder       bool                 `json:"is_folder"`
	ParentFolderID *primitive.ObjectID  `json:"parent_folder_id"`
	CreatedAt      time.Time            `json:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at"`
	OpenedAt       *time.Time           `json:"opened_at"`
	Description    string               `json:"description"`
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
	UserID             primitive.ObjectID       `json:"user_id"`
	FilePermissionType enums.FilePermissionType `json:"permission_type"`
	UserImage          string                   `json:"user_image"`
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
		TotalSize:      file.TotalSize,
		OwnerID:        file.OwnerID,
		IsFolder:       file.IsFolder,
		ParentFolderID: file.ParentFolderID,
		CreatedAt:      file.CreatedAt,
		UpdatedAt:      file.UpdatedAt,
		OpenedAt:       file.OpenedAt,
		Description:    file.Description,
		TagIDs:         file.TagIDs,
		StorageDetails: storageDetailOutput,
	}
}
