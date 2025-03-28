package impl

import (
	"context"
	"sync"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/constant"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/exception"
	commonModel "github.com/baothaihcmut/Bibox/storage-app/internal/common/models"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/response"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/models"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/presenters"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (f *FileInteractorImpl) RecoverFile(ctx context.Context, input *presenters.RecoverFileInput) (*presenters.RecoverFileOutput, error) {
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
	file, err := f.checkOwnerOfFile(ctx, fileId, userId)
	if err != nil {
		return nil, err
	}
	if file == nil {
		return nil, exception.ErrFileNotFound
	}
	if !file.IsDeleted {
		return nil, exception.ErrFileIsNotInBin
	}
	file.Recover()
	file.ParentFolderID = input.DestinationFolderId
	updateFile := []*models.File{file}
	wg := sync.WaitGroup{}
	errCh := make(chan error, 1)
	updateCh := make(chan *models.File, 100)
	doneCh := make(chan struct{}, 1)
	//update parent
	go func() {
		for file := range updateCh {
			updateFile = append(updateFile, file)
		}
		doneCh <- struct{}{}
	}()
	if input.DestinationFolderId != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			parentFolder, err := f.fileRepo.FindFileById(ctx, *input.DestinationFolderId)
			if err != nil {
				errCh <- err
				return
			}
			if !parentFolder.IsFolder {
				errCh <- exception.ErrFileIsFolder
				return
			}
			parentFolder.IncrementTotalSize(file.TotalSize)
			updateCh <- parentFolder
			if parentFolder.ParentFolderID != nil {
				//find all parent of parent
				parentFolders, err := f.fileRepo.FindAllParentFolder(ctx, parentFolder.ID)
				if err != nil {
					errCh <- err
					return
				}
				for _, folder := range parentFolders {
					folder.IncrementTotalSize(file.TotalSize)
					updateCh <- folder
				}
			}
		}()
	}
	if file.IsFolder {
		wg.Add(1)
		go func() {
			defer wg.Done()
			subFiles, err := f.fileRepo.FindSubFileRecursive(ctx, file.ID)
			if err != nil {
				errCh <- err
			}
			for _, subFile := range subFiles {
				subFile.Recover()
				updateCh <- subFile
			}
		}()
	}
	go func() {
		wg.Wait()
		close(updateCh)
	}()

	select {
	case err := <-errCh:
		return nil, err
	case <-doneCh:
	}
	session, err := f.mongoService.BeginTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer f.mongoService.EndTransansaction(ctx, session)
	if err := f.fileRepo.BulkUpdateFile(ctx, updateFile); err != nil {
		return nil, f.mongoService.RollbackTransaction(ctx, session)
	}
	if err := f.mongoService.CommitTransaction(ctx, session); err != nil {
		return nil, err
	}
	f.logger.Info(ctx, map[string]any{
		"file_id": fileId,
	}, "Recover file success")
	return &presenters.RecoverFileOutput{
		Files: lo.Map(updateFile, func(item *models.File, _ int) *response.FileOutput {
			return response.MapFileToFileOutput(item)
		}),
	}, nil
}
