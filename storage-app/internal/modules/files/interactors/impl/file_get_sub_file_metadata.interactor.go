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

func (f *FileInteractorImpl) GetSubFileMetaData(ctx context.Context, input *presenters.GetSubFileMetaDataInput) (*presenters.GetSubFileMetaDataOutput, error) {
	//check permission
	fileId, err := primitive.ObjectIDFromHex(input.FileId)
	if err != nil {
		return nil, exception.ErrInvalidObjectId
	}
	userCtx, _ := ctx.Value(constant.UserContext).(*commonModel.UserContext)
	userId, _ := primitive.ObjectIDFromHex(userCtx.Id)
	//check permssion
	file, err := f.checkFilePermission(ctx, fileId, userId, enums.ViewPermission)
	if err != nil {
		return nil, err
	}
	subFiles, err := f.fileRepo.FindFileByParentFolderIdAndIsDeleted(ctx, file.ID, input.IsDeleted)
	if err != nil {
		return nil, err
	}
	return &presenters.GetSubFileMetaDataOutput{
		SubFiles: lo.Map(subFiles, func(item *models.File, _ int) *response.FileOutput {
			return response.MapFileToFileOutput(item)

		}),
	}, nil
}
