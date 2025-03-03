package presenters

import (
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/tags/presenters"
)

type GetFileTagsInput struct {
	Id string `uri:"id" binding:"required"`
}

type GetFileTagsOutput struct {
	Tags []*presenters.TagOutput `json:"tags"`
}
