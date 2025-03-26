package presenters

import (
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/response"
)

type GetFileTagsInput struct {
	Id string `uri:"id" validate:"required"`
}

type GetFileTagsOutput struct {
	Tags []*response.TagOutput `json:"tags"`
}
