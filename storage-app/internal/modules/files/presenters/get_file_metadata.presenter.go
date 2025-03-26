package presenters

import "github.com/baothaihcmut/Bibox/storage-app/internal/common/response"

type GetFileMetaDataInput struct {
	Id string `uri:"id" validate:"required"`
}

type GetFileMetaDataOuput struct {
	*response.FileOutput
}
