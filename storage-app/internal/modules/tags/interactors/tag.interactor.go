package interactors

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/tags/presenters"
)

type TagInteractor interface {
	GetAllTags(ctx context.Context, input *presenters.SerchTagsInput) (*presenters.SearchTagsOutput, error)
}
