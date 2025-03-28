package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/constant"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/exception"
	commonModel "github.com/baothaihcmut/Bibox/storage-app/internal/common/models"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/response"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/models"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/presenters"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (f *FileInteractorImpl) SoftDeleteFile(ctx context.Context, input *presenters.SoftDeleteFileInput) (*presenters.SoftDeleteFileOuput, error) {
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
	file, err := f.checkFilePermission(ctx, fileId, userId, enums.EditPermission)
	if err != nil {
		return nil, err
	}
	//decrease all parent folder size
	updateFile := make([]*models.File, 0)
	if file.ParentFolderID != nil && file.TotalSize > 0 {
		parentFolders, err := f.fileRepo.FindAllParentFolder(ctx, file.ID)
		if err != nil {
			return nil, err
		}
		for _, parentFolder := range parentFolders {
			parentFolder.DecreaseTotalSize(file.TotalSize)
			updateFile = append(updateFile, parentFolder)
		}
	}
	//remove parent
	file.RemoveParent()
	//soft delete
	file.Delete()
	updateFile = append(updateFile, file)
	//find all subfile if file is folder
	if file.IsFolder {
		subFiles, err := f.fileRepo.FindSubFileRecursive(ctx, file.ID)
		if err != nil {
			return nil, err
		}
		for _, subFile := range subFiles {
			subFile.Delete()
			updateFile = append(updateFile, subFile)
		}
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
	}, "Soft delete file success")
	return &presenters.SoftDeleteFileOuput{
		Files: lo.Map(updateFile, func(item *models.File, _ int) *response.FileOutput {
			return response.MapFileToFileOutput(item)
		}),
	}, nil
}
