package presenters

import "github.com/baothaihcmut/Bibox/storage-app/internal/common/response"

type GetSubFileMetaDataInput struct {
	FileId string `uri:"id"`
}

type GetSubFileMetaDataOutput struct {
	SubFiles []*response.FileOutput `json:"sub_files"`
}
