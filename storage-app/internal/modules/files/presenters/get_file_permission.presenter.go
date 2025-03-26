package presenters

import "github.com/baothaihcmut/Bibox/storage-app/internal/common/response"

type FilePermissionUserInfo struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Image     string `json:"image"`
}

type GetFilePermissionInput struct {
	Id string `uri:"id" validate:"required"`
}

type FilePermssionWithUserOutput struct {
	*response.FilePermissionOuput

	User *FilePermissionUserInfo `json:"user"`
}

type GetFilePermissionOuput struct {
	Permissions []*FilePermssionWithUserOutput `json:"permissions"`
}
