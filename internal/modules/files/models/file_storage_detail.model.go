package models

type FileStorageDetail struct {
	Size            int    `bson:"size"`
	FileType        string `bson:"file_type"`
	IsUploaded      bool   `bson:"is_uploaded"`
	StorageProvider string `bson:"storage_provider"`
	StorageKey      string `bson:"storage_key"`
	StorageBucket   string `bson:"storage_bucket"`
}

func NewFileStorageDetail(
	size int,
	fileType string,
	isUploaded bool,
	storageProvider string,
	storageKey string,
	storageBucket string,
) *FileStorageDetail {
	return &FileStorageDetail{
		Size:            size,
		FileType:        fileType,
		IsUploaded:      isUploaded,
		StorageKey:      storageKey,
		StorageProvider: storageProvider,
		StorageBucket:   storageBucket,
	}
}
