package presenters

import (
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/response"
)

type GetFileTagsInput struct {
	Id string `uri:"id" binding:"required"`
}

type GetFileTagsOutput struct {
	Tags []*response.TagOutput `json:"tags"`
}
