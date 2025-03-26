package presenters

import "github.com/baothaihcmut/Bibox/storage-app/internal/common/response"

type UploadFolderInput struct {
	Data     *CreateFileInput     `json:"data"`
	SubFiles []*UploadFolderInput `json:"sub_files"`
}

type FileWithPathOutput struct {
	*response.FileOutput
	Path         string `json:"path"`
	PutObjectUrl string `json:"put_object_url"`
	UrlExpiry    int    `json:"url_expiry"`
}

type UploadFolderOutput struct {
	Files []*FileWithPathOutput `json:"files"`
}
