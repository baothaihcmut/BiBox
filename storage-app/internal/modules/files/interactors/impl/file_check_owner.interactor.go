package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/exception"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (f *FileInteractorImpl) checkOwnerOfFile(ctx context.Context, fileId primitive.ObjectID, userId primitive.ObjectID) (*models.File, error) {
	file, err := f.fileRepo.FindFileById(ctx, fileId, false)
	if err != nil {
		return nil, err
	}
	if file == nil {
		return nil, exception.ErrFileNotFound
	}
	if file.OwnerID != userId {
		return nil, exception.ErrPermissionDenied
	}
	return file, nil
}
