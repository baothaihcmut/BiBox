package presenters

type HardDeleteFileInput struct {
	Id string `uri:"id" validate:"required"`
}

type HardDeleteFileOutput struct{}
