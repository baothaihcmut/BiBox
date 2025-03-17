package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/constant"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
	commonModel "github.com/baothaihcmut/Bibox/storage-app/internal/common/models"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/presenters"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetFileMetaData implements FileInteractor.
func (f *FileInteractorImpl) GetFileMetaData(ctx context.Context, input *presenters.GetFileMetaDataInput) (*presenters.GetFileMetaDataOuput, error) {
	userCtx := ctx.Value(constant.UserContext).(*commonModel.UserContext)
	userId, _ := primitive.ObjectIDFromHex(userCtx.Id)
	fileId, _ := primitive.ObjectIDFromHex(input.Id)

	//check permission
	file, err := f.checkFilePermission(ctx, fileId, userId, enums.ViewPermission)
	if err != nil {
		return nil, err
	}
	return &presenters.GetFileMetaDataOuput{
		FileOutput: presenters.MapFileToFileOutput(file),
	}, nil
}
