package presenters

import "github.com/baothaihcmut/Bibox/storage-app/internal/common/response"

type SerchTagsInput struct {
	Query  string `form:"query"`
	Offset int    `form:"offset"`
	Limit  int    `form:"limit"`
}

type SearchTagsOutput struct {
	Data       []*response.TagOutput       `json:"tags"`
	Pagination response.PaginationResponse `json:"pagination"`
}
