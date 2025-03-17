package impl

import (
	"context"
	"time"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/constant"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/exception"
	commonModel "github.com/baothaihcmut/Bibox/storage-app/internal/common/models"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/storage"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/presenters"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetFileDownloadUrl implements FileInteractor.
func (f *FileInteractorImpl) GetFileDownloadUrl(ctx context.Context, input *presenters.GetFileDownloadUrlInput) (*presenters.GetFileDownloadUrlOutput, error) {
	//check permission and file find
	userCtx := ctx.Value(constant.UserContext).(*commonModel.UserContext)
	userId, _ := primitive.ObjectIDFromHex(userCtx.Id)
	fileId, _ := primitive.ObjectIDFromHex(input.Id)
	file, err := f.checkFilePermission(ctx, fileId, userId, enums.ViewPermission)
	if err != nil {
		return nil, err
	}
	//check if file is folder
	if file.IsFolder {
		return nil, exception.ErrFileIsFolder
	}
	//get url
	url, err := f.storageService.GetPresignUrl(ctx, storage.GetPresignUrlArg{
		Method:      storage.PresignUrlGetMethod,
		Key:         file.StorageDetail.StorageKey,
		ContentType: file.StorageDetail.MimeType,
		Expiry:      3 * time.Hour,
		Preview:     input.Preview,
	})
	if err != nil {
		return nil, err
	}
	return &presenters.GetFileDownloadUrlOutput{
		FileName:    file.Name,
		Url:         url,
		Expiry:      1,
		Method:      "GET",
		ContentType: file.StorageDetail.MimeType,
	}, nil
}
