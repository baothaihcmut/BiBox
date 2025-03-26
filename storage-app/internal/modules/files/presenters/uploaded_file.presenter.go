package presenters

import "github.com/baothaihcmut/Bibox/storage-app/internal/common/response"

type UploadedFileInput struct {
	Id string `uri:"id" bind:"required"`
}

type UploadedFileOutput struct {
	*response.FileOutput
}
