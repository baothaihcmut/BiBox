package presenters

type DeleteFilePermissionOfUserInput struct {
	FileId string `uri:"id" validate:"required"`
	UserId string `uri:"userId" validate:"required"`
}

type DeleteFilePermissionOfUserOutput struct{}
