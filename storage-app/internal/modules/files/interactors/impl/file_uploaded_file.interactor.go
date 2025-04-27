package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/exception"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/lock"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/response"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/presenters"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetFilePermissions implements FileInteractor.
func (f *FileInteractorImpl) UploadedFile(ctx context.Context, input *presenters.UploadedFileInput) (*presenters.UploadedFileOutput, error) {
	//check file exist
	fileId, err := primitive.ObjectIDFromHex(input.Id)
	if err != nil {
		return nil, exception.ErrInvalidObjectId
	}
	file, err := f.fileRepo.FindFileById(ctx, fileId)
	if err != nil {
		return nil, err
	}
	if file == nil {
		return nil, exception.ErrFileNotFound
	}
	if file.IsFolder {
		return nil, exception.ErrFileIsFolder
	}
	if file.StorageDetail.IsUploading {
		return nil, exception.ErrFileIsNotUploading
	}
	file.StorageDetail.IsUploaded = true
	file.StorageDetail.IsUploading = false

	//update db
	err = f.fileRepo.UpdateFile(ctx, file)
	if err != nil {
		return nil, err
	}
	//release lock
	if err := f.distrutedLockService.ReleaseLock(ctx, lockKey+file.ID.Hex(), input.UploadLockValue); err != nil {
		if err == lock.ErrKeyNotFound {
			f.logger.Errorf(ctx, map[string]any{
				"lock_key": lockKey + fileId.Hex(),
			}, "Lock key not found")
			return nil, exception.ErrFileIsNotUploading
		}
		if err == lock.ErrValueMisMatch {
			return nil, exception.ErrUploadFileLockValueMismatch
		}
		return nil, err
	}
	return &presenters.UploadedFileOutput{
		FileOutput: response.MapFileToFileOutput(file),
	}, nil
}
