package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/response"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/users/models"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/users/presenters"
	"github.com/samber/lo"
)

func (u *UserInteractorImpl) SearchUserByEmail(ctx context.Context, input *presenters.SearchUserInput) (*presenters.SearchUserOuput, error) {
	users, count, err := u.userRepo.FindUserRegexAndCount(ctx, input.Query, input.Limit, input.Offset)
	if err != nil {
		return nil, err
	}
	return &presenters.SearchUserOuput{
		Data: lo.Map(users, func(item *models.User, _ int) *response.UserOutput {
			return response.MapToUserOutput(item)

		}),
		Pagination: response.InitPaginationResponse(count, *input.Limit, *input.Offset),
	}, nil
}
