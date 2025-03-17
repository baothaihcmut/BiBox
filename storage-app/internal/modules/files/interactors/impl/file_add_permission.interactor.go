package impl

import (
	"context"
	"sync"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/constant"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/exception"
	commonModel "github.com/baothaihcmut/Bibox/storage-app/internal/common/models"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/presenters"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"

	permissionModel "github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/models"

	permissionPresenter "github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/presenters"
)

func (f *FileInteractorImpl) AddFilePermission(ctx context.Context, input *presenters.AddFilePermissionInput) (*presenters.AddFilePermissionOutput, error) {
	fileId, err := primitive.ObjectIDFromHex(input.FileId)
	if err != nil {
		return nil, exception.ErrInvalidObjectId
	}
	userCtx, _ := ctx.Value(constant.UserContext).(*commonModel.UserContext)
	userId, _ := primitive.ObjectIDFromHex(userCtx.Id)
	//check permssion
	file, err := f.checkFilePermission(ctx, fileId, userId, enums.EditPermission)
	if err != nil {
		return nil, err
	}
	//handle add permission to file
	if !file.IsFolder {

		permissions := make([]*permissionModel.FilePermission, len(input.Permissions))
		for idx, permission := range input.Permissions {
			if permission.PermissionType == nil {

			}
			permissions[idx] = permissionModel.NewFilePermission(
				fileId,
				permission.UserId,
				enums.FilePermissionType(*permission.PermissionType),
				true,
				nil,
			)
		}
		//add to db
		err = f.filePermissionRepo.BulkCreatePermission(ctx, permissions)
		if err != nil {
			return nil, err
		}
		return &presenters.AddFilePermissionOutput{
			Permissions: lo.Map(permissions, func(item *permissionModel.FilePermission, _ int) *permissionPresenter.FilePermissionOuput {
				return permissionPresenter.MapToOuput(item)
			}),
		}, nil
	}
	//get sub file
	subFiles, err := f.fileRepo.FindSubFileRecursive(ctx, file.ID)
	if err != nil {
		return nil, err
	}
	//build file tree structure
	fileStructure := f.fileStructureService.BuildFileStructureTree(ctx, append(subFiles, file))
	permissionCh := make(chan []*permissionModel.FilePermission, len(input.Permissions))
	wg := sync.WaitGroup{}
	for _, permission := range input.Permissions {
		wg.Add(1)
		go func() {
			defer wg.Done()
			permissions := f.fileStructureService.ExtractPermissionFromFileStructure(
				ctx,
				permission.UserId,
				fileStructure,
				permission.AdditionPermission)
			permissionCh <- permissions
		}()
	}
	go func() {
		wg.Wait()
		close(permissionCh)
	}()
	targetPermissions := make([]*permissionModel.FilePermission, 0)
	for permissions := range permissionCh {
		targetPermissions = append(targetPermissions, permissions...)
	}
	//insert to db
	err = f.filePermissionRepo.BulkCreatePermission(ctx, targetPermissions)
	if err != nil {
		return nil, err
	}
	return &presenters.AddFilePermissionOutput{
		Permissions: lo.Map(targetPermissions, func(item *permissionModel.FilePermission, _ int) *permissionPresenter.FilePermissionOuput {
			return permissionPresenter.MapToOuput(item)
		}),
	}, nil
}
