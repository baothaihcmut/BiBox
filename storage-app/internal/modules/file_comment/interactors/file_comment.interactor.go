package interactors

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_comment/presenters"
)

type FileCommentInteractor interface {
	CreateFileComment(ctx context.Context, input *presenters.CreateFileCommentInput) (*presenters.CreateFileCommentOutput, error)
}
