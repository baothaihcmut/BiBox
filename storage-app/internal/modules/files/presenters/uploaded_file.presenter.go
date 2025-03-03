package presenters

type UploadedFileInput struct {
	Id string `uri:"id" bind:"required"`
}

type UploadedFileOutput struct {
	*FileOutput
}
