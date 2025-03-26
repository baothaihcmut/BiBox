package presenters

import (
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/response"
)

type FindFileOfUserInput struct {
	IsFolder *bool   `form:"is_folder"`
	FileType *string `form:"mime_type"`
	SortBy   string  `form:"sort_by" validate:"required"`
	IsAsc    bool    `form:"is_asc" validate:"required"`
	Offset   int     `form:"offset" validate:"required"`
	Limit    int     `form:"limit" validate:"required"`
}

type FindFileOfUserOuput struct {
	Data       []*response.FileWithPermissionOutput `json:"data"`
	Pagination response.PaginationResponse          `json:"pagination"`
}
