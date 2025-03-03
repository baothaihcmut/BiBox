package presenters

import "github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/presenters"

type FilePermissionUserInfo struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Image     string `json:"image"`
}

type GetFilePermissionInput struct {
	Id string `uri:"id" binding:"required"`
}

type FilePermssionWithUserOutput struct {
	*presenters.FilePermissionOuput
	User *FilePermissionUserInfo `json:"user"`
}

type GetFilePermissionOuput struct {
	Permissions []*FilePermssionWithUserOutput `json:"permissions"`
}
