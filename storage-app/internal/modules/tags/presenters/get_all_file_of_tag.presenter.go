package presenters

import "github.com/baothaihcmut/Bibox/storage-app/internal/common/response"

type GetAllFilOfTagInput struct {
	Id     string `uri:"id" validate:"required"`
	Limit  int    `form:"limit" validate:"required"`
	Offset int    `form:"offset"`
	SortBy string `form:"sort_by" validate:"required" `
	IsAsc  bool   `form:"is_asc"`
}

type GetAllFileOfTagOutput struct {
	Data       []*response.FileWithPermissionOutput `json:"data"`
	Pagination response.PaginationResponse          `json:"pagination"`
}
