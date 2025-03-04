package presenters

type GetFileMetaDataInput struct {
	Id string `uri:"id" binding:"required"`
}

type GetFileMetaDataOuput struct {
	*FileOutput
}
