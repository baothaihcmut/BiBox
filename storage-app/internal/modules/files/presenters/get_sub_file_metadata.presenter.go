package presenters

type GetSubFileMetaDataInput struct {
	FileId string `uri:"id"`
}

type GetSubFileMetaDataOutput struct {
	SubFiles []*FileOutput `json:"sub_files"`
}
