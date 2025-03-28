package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/constant"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
	commonModel "github.com/baothaihcmut/Bibox/storage-app/internal/common/models"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/response"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/models"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/presenters"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (f *FileInteractorImpl) GetFilePermissions(ctx context.Context, input *presenters.GetFilePermissionInput) (*presenters.GetFilePermissionOuput, error) {
	userCtx := ctx.Value(constant.UserContext).(*commonModel.UserContext)
	userId, _ := primitive.ObjectIDFromHex(userCtx.Id)
	fileId, _ := primitive.ObjectIDFromHex(input.Id)
	file, err := f.checkFilePermission(ctx, fileId, userId, enums.ViewPermission)
	if err != nil {
		return nil, err
	}
	permission, err := f.filePermissionRepo.FindFilePermissionWithUser(ctx, file.ID)
	if err != nil {
		return nil, err
	}
	return &presenters.GetFilePermissionOuput{
		Permissions: lo.Map(permission, func(item *models.FilePermissionWithUser, _ int) *presenters.FilePermssionWithUserOutput {
			return &presenters.FilePermssionWithUserOutput{
				FilePermissionOuput: response.MapToFilePermissionOutput(item.FilePermission),

				User: &presenters.FilePermissionUserInfo{
					Email:     item.User.Email,
					FirstName: item.User.FirstName,
					LastName:  item.User.LastName,
					Image:     item.User.Image,
				},
			}
		}),
	}, nil
}
