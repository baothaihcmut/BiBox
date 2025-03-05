package presenters

import (
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/response"
)

type FindFileOfUserInput struct {
	ParentFolderId *string `form:"parent_folder_id"`
	IsFolder       *bool   `form:"is_folder"`
	FileType       *string `form:"mime_type"`
	SortBy         string  `form:"sort_by" bind:"required"`
	IsAsc          bool    `form:"is_asc" bind:"required"`
	Offset         int     `form:"offset" bind:"required"`
	Limit          int     `form:"limit" bind:"required"`
}

type FileWithPermissionOutput struct {
	*FileOutput
	Permissions []*PermissionOfFileOuput `json:"permissions"`
}

type FindFileOfUserOuput struct {
	Data       []*FileWithPermissionOutput `json:"data"`
	Pagination response.PaginationResponse `json:"pagination"`
}
