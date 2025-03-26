package impl

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/constant"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
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

func (f *FileInteractorImpl) AddFilePermission(ctx context.Context, input *presenters.AddFilePermissionInput) (*presenters.AddFilePermissionOutput, error) {
	//check user is owner of file
	userContext := ctx.Value(constant.UserContext).(*commonModel.UserContext)
	//tranform id
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
	wg := sync.WaitGroup{}
	errCh := make(chan error, 1)
	doneCh := make(chan struct{}, 1)
	var allFiles []*models.File
	//check user exist
	wg.Add(1)
	go func() {
		defer wg.Done()
		users, err := f.userRepo.FindUserIdIds(ctx, lo.Map(input.UserPermissions, func(item presenters.UserPermission, _ int) primitive.ObjectID {
			return item.UserId
		}))
		if err != nil {
			errCh <- err
			return
		}
		//init set of user id input
		mapUserId := make(map[primitive.ObjectID]struct{})
		for _, userPermission := range input.UserPermissions {
			mapUserId[userPermission.UserId] = struct{}{}
		}
		for _, user := range users {
			if _, exist := mapUserId[user.ID]; exist {
				delete(mapUserId, user.ID)
			}
		}
		if len(mapUserId) > 0 {
			errCh <- exception.ErrUserNotFound
			return
		}
	}()
	//find sub file
	if file.IsFolder {
		wg.Add(1)
		go func() {
			defer wg.Done()
			allFiles, err = f.fileRepo.FindSubFileRecursive(ctx, file.ID)
			if err != nil {
				f.logger.Errorf(ctx, nil, "Error find sub file of folder: ", err)
				errCh <- err
				return
			}
			//append parent
			allFiles = append(allFiles, file)
		}()
	} else {
		allFiles = []*models.File{file}
	}
	//wait for done check phase
	go func() {
		wg.Wait()
		doneCh <- struct{}{}
	}()
	select {
	case err = <-errCh:
		return nil, err
	case <-doneCh:
	}
	permissionIds := make([]repositories.FilePermissionId, 0, len(input.UserPermissions)*len(allFiles))
	mapFilePermissionType := make(map[string]struct {
		PermissionType enums.FilePermissionType
		ExpireAt       *time.Time
	})
	for _, userPermission := range input.UserPermissions {
		for _, file := range allFiles {
			permissionIds = append(permissionIds, repositories.FilePermissionId{
				FileId: file.ID,
				UserId: userPermission.UserId,
			})
			key := file.ID.Hex() + "_" + userPermission.UserId.Hex()
			mapFilePermissionType[key] = struct {
				PermissionType enums.FilePermissionType
				ExpireAt       *time.Time
			}{
				PermissionType: userPermission.PermissionsType,
				ExpireAt:       userPermission.ExpireAt,
			}
		}
	}
	//file permission exist
	filePermissionsExist, err := f.filePermissionRepo.FindPermissionByIds(ctx, permissionIds)
	if err != nil {
		f.logger.Errorf(ctx, nil, "Error find permission by id list: ", err)
		return nil, err
	}
	updatePermissions := make([]*permissionModel.FilePermission, 0)
	createPermissions := make([]*permissionModel.FilePermission, 0)
	for _, filePermissionOld := range filePermissionsExist {
		key := filePermissionOld.FileID.Hex() + "_" + filePermissionOld.UserID.Hex()
		if newPermission, exist := mapFilePermissionType[key]; exist {
			if filePermissionOld.FilePermissionType < newPermission.PermissionType {
				filePermissionOld.FilePermissionType = newPermission.PermissionType
				filePermissionOld.ExpireAt = newPermission.ExpireAt
				updatePermissions = append(updatePermissions, filePermissionOld)
			}
			delete(mapFilePermissionType, key)
		}
	}
	for k, v := range mapFilePermissionType {
		splitKey := strings.Split(k, "_")
		fileId, _ := primitive.ObjectIDFromHex(splitKey[0])
		userId, _ := primitive.ObjectIDFromHex(splitKey[1])
		createPermissions = append(createPermissions, permissionModel.NewFilePermission(
			fileId,
			userId,
			v.PermissionType,
			true,
			v.ExpireAt,
		))
	}
	//save to db
	wgSave := sync.WaitGroup{}
	errSave := make(chan error, 1)
	doneSave := make(chan struct{}, 1)
	session, err := f.mongoService.BeginTransaction(ctx)
	if err != nil {
		f.logger.Errorf(ctx, nil, "Error init traction: ", err)
		return nil, err
	}
	defer f.mongoService.EndTransansaction(ctx, session)
	if len(updatePermissions) > 0 {
		wgSave.Add(1)
		go func() {
			defer wgSave.Done()
			err = f.filePermissionRepo.BulkUpdatePermission(ctx, updatePermissions)
			if err != nil {
				errSave <- f.mongoService.RollbackTransaction(ctx, session)
			}
		}()
	}
	if len(createPermissions) > 0 {
		wgSave.Add(1)
		go func() {
			defer wgSave.Done()
			err = f.filePermissionRepo.BulkCreatePermission(ctx, createPermissions)
			if err != nil {
				errSave <- f.mongoService.RollbackTransaction(ctx, session)
			}

		}()
	}
	go func() {
		wgSave.Wait()
		doneSave <- struct{}{}
	}()
	select {
	case err := <-errCh:
		return nil, err
	case <-doneSave:
	}
	if err := f.mongoService.CommitTransaction(ctx, session); err != nil {
		f.logger.Errorf(ctx, nil, "Error commit transaction: ", err)
		return nil, err
	}
	return &presenters.AddFilePermissionOutput{
		Permissions: lo.Map(append(createPermissions, updatePermissions...), func(item *permissionModel.FilePermission, _ int) *response.FilePermissionOuput {
			return response.MapToFilePermissionOutput(item)
		}),
	}, nil

}
