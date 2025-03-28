package presenters

import "github.com/baothaihcmut/Bibox/storage-app/internal/common/response"

type SoftDeleteFileInput struct {
	Id string `uri:"id" validate:"required"`
}

type SoftDeleteFileOuput struct {
	files []*response.FileOutput
}
