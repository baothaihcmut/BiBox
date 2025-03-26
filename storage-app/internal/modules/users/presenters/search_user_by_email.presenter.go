package presenters

import "github.com/baothaihcmut/Bibox/storage-app/internal/common/response"

type SearchUserByEmailInput struct {
	Email  string `form:"email" binding:"required"`
	Offset *int   `form:"offset" binding:"gte=0"`
	Limit  *int   `form:"limit" binding:"gte=0"`
}

type SearchUserByEmailOuput struct {
	Data       []*UserOutput               `json:"data"`
	Pagination response.PaginationResponse `json:"pagination"`
}
