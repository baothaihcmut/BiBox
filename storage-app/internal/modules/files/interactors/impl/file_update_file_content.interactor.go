package impl

import (
	"context"
	"sync"
	"time"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/exception"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/response"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/storage"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/utils"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/models"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/presenters"
	userModel "github.com/baothaihcmut/Bibox/storage-app/internal/modules/users/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (f *FileInteractorImpl) UpdateFileContent(ctx context.Context, input *presenters.UpdateFileContentInput) (*presenters.UpdateFileContentOutput, error) {
	userContext := utils.GetUserContext(ctx)
	lockFileKey := lockKey + input.Id
	lockFileValue := uuid.New().String()
	ok, err := f.distrutedLockService.AcquireLock(ctx, lockFileKey, lockFileValue, 3, time.Second*3, time.Minute*30)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, exception.ErrFileIsUploading
	}

	fileId, err := primitive.ObjectIDFromHex(input.Id)
	if err != nil {
		return nil, exception.ErrInvalidObjectId
	}
	userId, _ := primitive.ObjectIDFromHex(userContext.Id)
	errCh := make(chan error, 1)
	userCh := make(chan *userModel.User, 1)
	fileCh := make(chan *models.File, 1)
	wgCheck := sync.WaitGroup{}
	wgCheck.Add(1)
	go func() {
		defer wgCheck.Done()
		file, err := f.checkFilePermission(ctx, fileId, userId, enums.EditPermission)
		if err != nil {
			errCh <- err
			return
		}
		if file.IsFolder {
			errCh <- exception.ErrFileIsFolder
			return
		}
		if file.StorageDetail.MimeType != enums.MapToMimeType(input.StorageDetail.MimeType, "application/pdf") {
			errCh <- exception.ErrMimeTypeMismatch
			return
		}
		fileCh <- file
	}()
	wgCheck.Add(1)
	go func() {
		defer wgCheck.Done()
		user, err := f.userRepo.FindUserById(ctx, userId)
		if err != nil {
			errCh <- err
			return
		}
		file := <-fileCh
		if err := user.DecreStorageSize(file.TotalSize); err != nil {
			errCh <- err
			return
		}
		if err := user.IncreStorageSize(input.StorageDetail.Size); err != nil {
			errCh <- err
			return
		}
		fileCh <- file
		userCh <- user
	}()
	wgCheck.Wait()
	select {
	case <-ctx.Done():
		return nil, nil
	case err = <-errCh:
		return nil, err
	default:
	}

	file := <-fileCh
	file.UpdatedAt = time.Now()
	file.StorageDetail.IsUploaded = false
	file.StorageDetail.IsUploading = true
	file.StorageDetail.Size = input.StorageDetail.Size
	wgSave := sync.WaitGroup{}
	errSave := make(chan error, 1)
	session, err := f.mongoService.BeginTransaction(ctx)
	if err != nil {
		f.logger.Errorf(ctx, nil, "Error init transaction mongo:", err)
		return nil, err
	}
	defer func() {
		if err != nil {
			f.mongoService.RollbackTransaction(ctx, session)
		}
		f.mongoService.EndTransansaction(ctx, session)
	}()
	wgSave.Add(1)
	go func() {
		defer wgSave.Done()

		if err := f.fileRepo.UpdateFile(ctx, file); err != nil {
			errSave <- err
			return
		}
	}()
	wgSave.Add(1)
	go func() {
		defer wgSave.Done()
		user := <-userCh
		if err := f.userRepo.UpdateUserStorageSize(ctx, user); err != nil {
			errSave <- err
		}
	}()
	wgSave.Wait()
	select {
	case <-ctx.Done():
		return nil, nil
	case err = <-errSave:
		return nil, err
	default:
	}
	if err := f.mongoService.CommitTransaction(ctx, session); err != nil {
		f.logger.Errorf(ctx, nil, "Error commit transaction: ", err)
		return nil, err
	}
	url, err := f.storageService.GetPresignUrl(ctx, storage.GetPresignUrlArg{
		Method:      storage.PresignUrlPutMethod,
		Key:         file.StorageDetail.StorageKey,
		ContentType: file.StorageDetail.MimeType,
		Expiry:      time.Minute * 30,
	})
	if err != nil {
		f.logger.Errorf(ctx, nil, "Error get presign url: ", err)
		return nil, err
	}
	return &presenters.UpdateFileContentOutput{
		FileOutput:      response.MapFileToFileOutput(file),
		PutObjectUrl:    url,
		UrlExpiry:       30,
		UploadLockValue: lockFileValue,
	}, nil

}
