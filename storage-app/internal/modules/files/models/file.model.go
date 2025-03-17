package models

import (
	"time"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileStorageDetailArg struct {
	Size            int
	MimeType        enums.MimeType
	StorageProvider string
	StorageBucket   string
}

type File struct {
	ID             primitive.ObjectID   `bson:"_id,omitempty"`
	Name           string               `bson:"name"`
	TotalSize      int                  `bson:"total_size"`
	OwnerID        primitive.ObjectID   `bson:"owner_id"`
	IsFolder       bool                 `bson:"is_folder"`
	ParentFolderID *primitive.ObjectID  `bson:"parent_folder_id"`
	CreatedAt      time.Time            `bson:"created_at"`
	UpdatedAt      time.Time            `bson:"updated_at"`
	OpenedAt       *time.Time           `bson:"opened_at"`
	Description    string               `bson:"description"`
	IsDeleted      bool                 `bson:"is_deleted"`
	DeletedAt      *time.Time           `bson:"deleted_at"`
	TagIDs         []primitive.ObjectID `bson:"tag_ids"`
	StorageDetail  *FileStorageDetail   `bson:"storage_detail"`
}

// Constructor for File with a random ObjectID
func NewFile(
	ownerID primitive.ObjectID,
	name string,
	parentFolderID *primitive.ObjectID,
	description string,
	isFolder bool,
	tags []primitive.ObjectID,
	storageDetail *FileStorageDetailArg) *File {
	id := primitive.NewObjectID()
	//key for storage
	key := uuid.New().String()
	//intit storage storage
	var storage *FileStorageDetail
	totalSize := 0
	if !isFolder {
		totalSize = storageDetail.Size
		storage = NewFileStorageDetail(
			storageDetail.Size,
			storageDetail.MimeType,
			false,
			storageDetail.StorageProvider,
			key, storageDetail.StorageBucket)
	}

	return &File{
		ID:             id, // Generate a random ObjectID
		OwnerID:        ownerID,
		Name:           name,
		TotalSize:      totalSize,
		IsFolder:       isFolder,
		ParentFolderID: parentFolderID,
		CreatedAt:      time.Now(), // Set CreatedAt to the current time
		UpdatedAt:      time.Now(), // Set UpdatedAt to the current time
		OpenedAt:       nil,        // Optionally set this to a default value
		Description:    description,
		IsDeleted:      false, // Default to not deleted
		DeletedAt:      nil,   // Optionally set this to a default value
		TagIDs:         tags,  // Pass the provided tags
		StorageDetail:  storage,
	}
}

func (f *File) IncrementTotalSize(size int) {
	f.TotalSize += size
}
