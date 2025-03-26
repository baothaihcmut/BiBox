package presenters

import "github.com/baothaihcmut/Bibox/storage-app/internal/common/response"

type GetSubFileOfFolderInput struct {
	Id       string  `uri:"id"`
	IsFolder *bool   `form:"is_folder"`
	FileType *string `form:"mime_type"`
	SortBy   string  `form:"sort_by" validate:"required"`
	IsAsc    bool    `form:"is_asc" validate:"required"`
	Offset   int     `form:"offset" validate:"required"`
	Limit    int     `form:"limit" validate:"required"`
}

type GetSubFileOfFolderOutput struct {
	Data       []*response.FileWithPermissionOutput `json:"data"`
	Pagination response.PaginationResponse          `json:"pagination"`
}
