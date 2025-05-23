package models

import "github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"

type FileStorageDetail struct {
	Size            int            `bson:"size"`
	MimeType        enums.MimeType `bson:"mime_type"`
	IsUploaded      bool           `bson:"is_uploaded"`
	IsUploading     bool           `bson:"is_uploading"`
	StorageProvider string         `bson:"storage_provider"`
	StorageKey      string         `bson:"storage_key"`
	StorageBucket   string         `bson:"storage_bucket"`
}

func NewFileStorageDetail(
	size int,
	fileType enums.MimeType,
	isUploaded bool,
	isUploading bool,
	storageProvider string,
	storageKey string,
	storageBucket string,
) *FileStorageDetail {
	return &FileStorageDetail{
		Size:            size,
		MimeType:        fileType,
		IsUploaded:      isUploaded,
		IsUploading:     isUploading,
		StorageKey:      storageKey,
		StorageProvider: storageProvider,
		StorageBucket:   storageBucket,
	}
}
