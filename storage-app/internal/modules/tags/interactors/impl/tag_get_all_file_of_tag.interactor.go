package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/constant"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/exception"
	commonModel "github.com/baothaihcmut/Bibox/storage-app/internal/common/models"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/response"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/models"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/repositories"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/tags/presenters"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (t *TagInteractorImpl) GetAllFileOfTag(ctx context.Context, input *presenters.GetAllFilOfTagInput) (*presenters.GetAllFileOfTagOutput, error) {
	userCtx := ctx.Value(constant.UserContext).(*commonModel.UserContext)
	userId, _ := primitive.ObjectIDFromHex(userCtx.Id)
	tagId, err := primitive.ObjectIDFromHex(input.Id)
	if err != nil {
		return nil, exception.ErrInvalidObjectId
	}
	sortBy := "created_at"
	if input.SortBy != "" {
		sortBy = input.SortBy
	}
	data, count, err := t.fileRepo.FindFileWithPermssionAndCount(
		ctx,
		repositories.FindFileWithPermissionArg{
			SortBy: sortBy,
			IsAsc:  input.IsAsc,
			Offset: input.Offset,
			Limit:  input.Limit,
			UserId: userId,
			TagId:  &tagId,
		},
	)

	if err != nil {
		return nil, err
	}
	return &presenters.GetAllFileOfTagOutput{
		Data: lo.Map(data, func(item *models.FileWithPermission, _ int) *response.FileWithPermissionOutput {
			filePermissions := make([]*response.PermissionOfFileOuput, 0, len(item.Permissions))
			for _, permission := range item.Permissions {
				filePermissions = append(filePermissions, &response.PermissionOfFileOuput{
					UserID:             permission.UserID,
					FilePermissionType: permission.FilePermissionType,
					UserImage:          permission.UserImage,
				})
			}
			return &response.FileWithPermissionOutput{
				FileOutput:  response.MapFileToFileOutput(&item.File),
				Permissions: filePermissions,
			}
		}),
		Pagination: response.InitPaginationResponse(
			int(count), input.Limit, input.Offset,
		),
	}, nil

}
