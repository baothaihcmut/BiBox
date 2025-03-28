package presenters

import (
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/response"
)

type FindFileOfUserInput struct {
	IsFolder *bool   `form:"is_folder"`
	FileType *string `form:"mime_type"`
	SortBy   string  `form:"sort_by" validate:"required"`
	IsAsc    bool    `form:"is_asc" validate:"required"`
	Offset   int     `form:"offset" validate:"gte=0"`
	Limit    int     `form:"limit" validate:"required,gt=0"`
}

type FindFileOfUserOuput struct {
	Data       []*response.FileWithPermissionOutput `json:"data"`
	Pagination response.PaginationResponse          `json:"pagination"`
}
