package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/constant"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/exception"
	commonModel "github.com/baothaihcmut/Bibox/storage-app/internal/common/models"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/response"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/repositories"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/presenters"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (f *FileInteractorImpl) GetFilePermissionOfUser(ctx context.Context, input *presenters.GetFilePermissionOfUserInput) (*presenters.GetFilePermissionOfUserOutput, error) {
	userCtx := ctx.Value(constant.UserContext).(*commonModel.UserContext)
	userId, _ := primitive.ObjectIDFromHex(userCtx.Id)
	fileId, err := primitive.ObjectIDFromHex(input.FileId)
	if err != nil {
		return nil, exception.ErrInvalidObjectId
	}
	//check file exist
	file, err := f.fileRepo.FindFileById(ctx, fileId)
	if err != nil {
		return nil, err
	}
	if file == nil {
		return nil, exception.ErrFileNotFound
	}

	permission, err := f.filePermissionRepo.FindFilePermissionById(ctx, repositories.FilePermissionId{
		FileId: fileId,
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}
	if permission == nil {
		return nil, exception.ErrPermissionDenied
	}
	return &presenters.GetFilePermissionOfUserOutput{
		FilePermissionOuput: *response.MapToFilePermissionOutput(permission),
	}, nil

}
