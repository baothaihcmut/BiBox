package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/constant"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/exception"
	commonModel "github.com/baothaihcmut/Bibox/storage-app/internal/common/models"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/response"
	permissionModel "github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/models"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/repositories"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/models"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/presenters"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (f *FileInteractorImpl) UpdateFilePermission(ctx context.Context, input *presenters.UpdateFilePermissionInput) (*presenters.UpdateFilePermissionOuput, error) {
	userContext := ctx.Value(constant.UserContext).(*commonModel.UserContext)

	userId, err := primitive.ObjectIDFromHex(userContext.Id)
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
	targetUserId, err := primitive.ObjectIDFromHex(input.UserId)
	if err != nil {
		return nil, exception.ErrInvalidObjectId
	}
	//get all subfile of file
	if !file.IsFolder {
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
		filePermission.FilePermissionType = input.PermissionType
		filePermission.ExpireAt = input.ExpireAt
		if err := f.filePermissionRepo.UpdatePermission(ctx, filePermission); err != nil {
			return nil, err
		}
		f.logger.Info(ctx, map[string]any{
			"file_id": file.ID,
			"user_id": targetUserId,
		}, "Update file permission success")
		return nil, err
	}
	//get all subfoler
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
	for _, filePermission := range filePermissions {
		filePermission.FilePermissionType = input.PermissionType
		filePermission.ExpireAt = input.ExpireAt
	}
	//save to db
	session, err := f.mongoService.BeginTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer f.mongoService.EndTransansaction(ctx, session)

	if err := f.filePermissionRepo.BulkUpdatePermission(ctx, filePermissions); err != nil {
		f.mongoService.RollbackTransaction(ctx, session)
		return nil, err
	}
	if err := f.mongoService.CommitTransaction(ctx, session); err != nil {
		f.logger.Errorf(ctx, nil, "Error commit transaction: ", err)
		return nil, err
	}
	return &presenters.UpdateFilePermissionOuput{
		Permissions: lo.Map(filePermissions, func(item *permissionModel.FilePermission, _ int) *response.FilePermissionOuput {
			return response.MapToFilePermissionOutput(item)
		}),
	}, nil
}
