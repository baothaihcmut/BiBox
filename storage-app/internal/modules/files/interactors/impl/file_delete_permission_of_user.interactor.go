package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/constant"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/exception"
	commonModel "github.com/baothaihcmut/Bibox/storage-app/internal/common/models"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/repositories"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/models"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/presenters"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (f *FileInteractorImpl) DeleteFilePermission(ctx context.Context, input *presenters.DeleteFilePermissionOfUserInput) (*presenters.DeleteFilePermissionOfUserOutput, error) {
	userContext := ctx.Value(constant.UserContext).(*commonModel.UserContext)
	//tranform id
	userId, err := primitive.ObjectIDFromHex(userContext.Id)
	if err != nil {
		return nil, exception.ErrInvalidObjectId
	}
	targetUserId, err := primitive.ObjectIDFromHex(input.UserId)
	if err != nil {
		return nil, exception.ErrInvalidObjectId
	}
	fileId, err := primitive.ObjectIDFromHex(input.FileId)
	if err != nil {
		return nil, exception.ErrInvalidObjectId
	}
	file, err := f.checkOwnerOfFile(ctx, fileId, userId)
	if err != nil {
		return nil, err
	}
	if !file.IsFolder {
		//check if permission exist
		filePermission, err := f.filePermissionRepo.FindFilePermissionById(ctx, repositories.FilePermissionId{
			FileId: file.ID,
			UserId: targetUserId,
		})
		if err != nil {
			return nil, err
		}
		if filePermission == nil {
			return nil, exception.ErrFilePermissionNotFound
		}
		if err := f.filePermissionRepo.DeletePermission(ctx, filePermission); err != nil {
			f.logger.Errorf(ctx, map[string]any{
				"file_id": file.ID,
				"user_id": targetUserId,
			}, "Error delete file permission of user: ", err)
			return nil, err
		}
		return &presenters.DeleteFilePermissionOfUserOutput{}, nil
	}

	//file all sub file of folder
	subFiles, err := f.fileRepo.FindSubFileRecursive(ctx, file.ID)
	if err != nil {
		return nil, err
	}
	filePermissionIds := lo.Map(append(subFiles, file), func(item *models.File, _ int) repositories.FilePermissionId {
		return repositories.FilePermissionId{
			FileId: item.ID,
			UserId: targetUserId,
		}
	})
	filePermissions, err := f.filePermissionRepo.FindPermissionByListId(ctx, filePermissionIds)
	if err != nil {
		return nil, err
	}
	//delete permission
	session, err := f.mongoService.BeginTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer f.mongoService.EndTransansaction(ctx, session)
	if err := f.filePermissionRepo.BulkDeletePermission(ctx, filePermissions); err != nil {
		return nil, f.mongoService.RollbackTransaction(ctx, session)
	}
	if err := f.mongoService.CommitTransaction(ctx, session); err != nil {
		return nil, err
	}
	return &presenters.DeleteFilePermissionOfUserOutput{}, nil
}
