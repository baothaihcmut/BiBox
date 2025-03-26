package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/response"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/tags/models"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/tags/presenters"
	"github.com/samber/lo"
)

func (t *TagInteractorImpl) GetAllTags(ctx context.Context, input *presenters.SerchTagsInput) (*presenters.SearchTagsOutput, error) {
	data, count, err := t.repo.FindAllTagsAndCount(ctx, input.Query, input.Limit, input.Offset)
	if err != nil {
		t.logger.Errorf(ctx, nil, "Error find all tags: %v", err)
		return nil, err
	}
	return &presenters.SearchTagsOutput{
		Data: lo.Map(data, func(item *models.Tag, _ int) *response.TagOutput {
			return response.MapToTagOuput(item)
		}),
		Pagination: response.InitPaginationResponse(
			count,
			input.Limit,
			input.Offset,
		),
	}, nil
}
