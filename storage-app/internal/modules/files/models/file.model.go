package models

import (
	"time"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
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
	OwnerID        primitive.ObjectID   `bson:"owner_id"`
	IsFolder       bool                 `bson:"is_folder"`
	ParentFolderID *primitive.ObjectID  `bson:"parent_folder_id"`
	CreatedAt      time.Time            `bson:"created_at"`
	UpdatedAt      time.Time            `bson:"updated_at"`
	OpenedAt       *time.Time           `bson:"opened_at"`
	HasPassword    bool                 `bson:"has_password"`
	Password       *string              `bson:"password"`
	Description    string               `bson:"description"`
	IsDeleted      bool                 `bson:"is_deleted"`
	DeletedAt      *time.Time           `bson:"deleted_at"`
	IsSecure       bool                 `bson:"is_secure"`
	TagIDs         []primitive.ObjectID `bson:"tag_ids"`
	StorageDetail  *FileStorageDetail   `bson:"storage_detail"`
}

// Constructor for File with a random ObjectID
func NewFile(
	ownerID primitive.ObjectID,
	name string,
	parentFolderID *primitive.ObjectID,
	description string,
	pasword *string,
	isFolder, hasPassword, isSecure bool,
	tags []primitive.ObjectID,
	storageDetail *FileStorageDetailArg) *File {
	id := primitive.NewObjectID()
	//key for storage
	key := time.Now().String() + id.Hex()
	//intit storage storage
	var storage *FileStorageDetail
	if !isFolder {
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
		IsFolder:       isFolder,
		ParentFolderID: parentFolderID,
		CreatedAt:      time.Now(), // Set CreatedAt to the current time
		UpdatedAt:      time.Now(), // Set UpdatedAt to the current time
		OpenedAt:       nil,        // Optionally set this to a default value
		HasPassword:    hasPassword,
		Password:       pasword, // Default to empty, can be set later
		Description:    description,
		IsDeleted:      false, // Default to not deleted
		DeletedAt:      nil,   // Optionally set this to a default value
		IsSecure:       isSecure,
		TagIDs:         tags, // Pass the provided tags
		StorageDetail:  storage,
	}
}
