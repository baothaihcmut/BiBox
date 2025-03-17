package impl

import (
	"context"
	"sync"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/exception"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (f *FileInteractorImpl) checkFilePermission(ctx context.Context, fileId primitive.ObjectID, userId primitive.ObjectID, permissionType enums.FilePermissionType) (*models.File, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	fileCh := make(chan *models.File, 1)
	errCh := make(chan error, 1)
	wg := sync.WaitGroup{}
	wg.Add(1)
	//check file exist
	go func() {
		defer wg.Done()
		file, err := f.fileRepo.FindFileById(ctx, fileId, false)
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
			}
			cancel()
			errCh <- err
			return
		}
		if file == nil {
			cancel()
			errCh <- exception.ErrFileNotFound
			return
		}
		fileCh <- file
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		hasPermission, err := f.filePermission.CheckPermission(ctx, fileId, userId, permissionType)
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
			}
			cancel()
			errCh <- err
			return
		}
		if !hasPermission {
			errCh <- exception.ErrPermissionDenied
			return
		}
	}()
	wg.Wait()
	select {
	case err := <-errCh:
		return nil, err
	default:
		return <-fileCh, nil
	}
}
