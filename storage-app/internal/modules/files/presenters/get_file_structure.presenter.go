package presenters

type GetFileStructureInput struct {
	Id string `uri:"id" binding:"required"`
}

type GetFileStructrueOuput struct {
	SubFiles []*FileOutput `json:"sub_files"`
}
