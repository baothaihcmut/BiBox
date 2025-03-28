package impl

import (
	"context"
	"sync"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/constant"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/exception"
	commonModel "github.com/baothaihcmut/Bibox/storage-app/internal/common/models"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/models"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/presenters"
	userModel "github.com/baothaihcmut/Bibox/storage-app/internal/modules/users/models"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (f *FileInteractorImpl) HardDeleteFile(ctx context.Context, input *presenters.HardDeleteFileInput) (*presenters.HardDeleteFileOutput, error) {
	userContext := ctx.Value(constant.UserContext).(*commonModel.UserContext)
	//tranform id
	userId, err := primitive.ObjectIDFromHex(userContext.Id)
	if err != nil {
		return nil, exception.ErrInvalidObjectId
	}

	fileId, err := primitive.ObjectIDFromHex(input.Id)
	if err != nil {
		return nil, exception.ErrInvalidObjectId
	}
	userCh := make(chan *userModel.User, 1)
	errCh := make(chan error, 1)
	totalSizeCh := make(chan int, 1)
	allFileIdsCh := make(chan []primitive.ObjectID, 1)
	allObjectId := make(chan []string, 1)
	doneCh := make(chan struct{}, 1)
	wg := sync.WaitGroup{}
	session, err := f.mongoService.BeginTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer f.mongoService.EndTransansaction(ctx, session)
	wg.Add(1)
	go func() {
		defer wg.Done()
		user, err := f.userRepo.FindUserById(ctx, userId)
		if err != nil {
			errCh <- err
			return
		}
		if user == nil {
			errCh <- exception.ErrUserNotFound
			return
		}
		userCh <- user
		//update file size
		if err := user.DecreStorageSize(<-totalSizeCh); err != nil {
			errCh <- err
			return
		}
		if err := f.userRepo.UpdateUserStorageSize(ctx, user); err != nil {
			errCh <- err
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		file, err := f.checkOwnerOfFile(ctx, fileId, userId)
		if err != nil {
			errCh <- err
			return
		}
		if file == nil {
			errCh <- exception.ErrFileNotFound
			return
		}
		totalSizeCh <- file.TotalSize
		deleteFile := []*models.File{file}
		if file.IsFolder {
			subFiles, err := f.fileRepo.FindSubFileRecursive(ctx, file.ID)
			if err != nil {
				errCh <- err
				return
			}
			deleteFile = append(deleteFile, subFiles...)
		}
		allFileIdsCh <- lo.Map(deleteFile, func(item *models.File, _ int) primitive.ObjectID {
			return item.ID
		})
		allObjectId <- lo.Map(lo.Filter(deleteFile, func(item *models.File, _ int) bool {
			return !item.IsFolder && item.StorageDetail != nil
		}), func(item *models.File, _ int) string {
			return item.StorageDetail.StorageKey
		})
		if err := f.fileRepo.BulkDeleteFile(ctx, deleteFile); err != nil {
			errCh <- err
			return
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := f.filePermissionRepo.DeletePermissionByListFileId(ctx, <-allFileIdsCh); err != nil {
			errCh <- err
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := f.storageService.BulkDelete(ctx, <-allObjectId); err != nil {
			errCh <- err
			return
		}
	}()

	go func() {
		wg.Wait()
		doneCh <- struct{}{}
	}()
	select {
	case err := <-errCh:
		f.mongoService.RollbackTransaction(ctx, session)
		return nil, err
	case <-doneCh:
	}
	if err := f.mongoService.CommitTransaction(ctx, session); err != nil {
		return nil, err
	}
	return &presenters.HardDeleteFileOutput{}, nil
}
