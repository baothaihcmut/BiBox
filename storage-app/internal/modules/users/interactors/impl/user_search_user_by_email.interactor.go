package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/response"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/users/models"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/users/presenters"
	"github.com/samber/lo"
)

func (u *UserInteractorImpl) SearchUserByEmail(ctx context.Context, input *presenters.SearchUserByEmailInput) (*presenters.SearchUserByEmailOuput, error) {
	users, count, err := u.userRepo.FindUserByEmailRegexAndCount(ctx, input.Email, input.Limit, input.Offset)
	if err != nil {
		return nil, err
	}
	return &presenters.SearchUserByEmailOuput{
		Data: lo.Map(users, func(item *models.User, _ int) *presenters.UserOutput {
			return presenters.MapToUserOutput(item)
		}),
		Pagination: response.InitPaginationResponse(count, *input.Limit, *input.Offset),
	}, nil
}
