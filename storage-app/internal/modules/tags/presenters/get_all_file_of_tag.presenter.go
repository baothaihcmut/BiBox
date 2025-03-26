package presenters

import "github.com/baothaihcmut/Bibox/storage-app/internal/common/response"

type GetAllFilOfTagInput struct {
	Id     string `uri:"id"`
	Limit  int    `form:"limit"`
	Offset int    `form:"offset"`
	SortBy string `form:"sort_by"`
	IsAsc  bool   `form:"is_asc"`
}

type GetAllFileOfTagOutput struct {
	Data       []*response.FileWithPermissionOutput `json:"data"`
	Pagination response.PaginationResponse          `json:"pagination"`
}
