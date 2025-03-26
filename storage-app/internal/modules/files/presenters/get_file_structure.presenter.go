package presenters

import "github.com/baothaihcmut/Bibox/storage-app/internal/common/response"

type GetFileStructureInput struct {
	Id string `uri:"id" binding:"required"`
}

type GetFileStructrueOuput struct {
	SubFiles []*response.FileOutput `json:"sub_files"`
}
