package presenters

import "github.com/baothaihcmut/Bibox/storage-app/internal/common/response"

type SearchUserInput struct {
	Query  string `form:"email"`
	Offset *int   `form:"offset" binding:"gte=0"`
	Limit  *int   `form:"limit" binding:"gte=0"`
}

type SearchUserOuput struct {
	Data       []*response.UserOutput      `json:"data"`
	Pagination response.PaginationResponse `json:"pagination"`
}
