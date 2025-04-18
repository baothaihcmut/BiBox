package presenters

import (
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/response"
)

type CreateFileCommentInput struct {
	FileId  string `uri:"id" validate:"required"`
	Content string `json:"content" validate:"required"`
}

type CreateFileCommentOutput struct {
	*response.FileCommentOutput
}
