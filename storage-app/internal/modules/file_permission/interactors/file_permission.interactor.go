package interactors

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/logger"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/models"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/presenters"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FilePermissionInteractor interface {
	GrantFilePermission(context.Context, *presenters.GrantFilePermissionInput) (*presenters.GrantFilePermissionOutput, error)
}

type FilePermissionInteractorImpl struct {
	filePermissionRepo repositories.FilePermissionRepository
	logger             logger.Logger
}

func NewFilePermissionInteractor(repo repositories.FilePermissionRepository, logger logger.Logger) FilePermissionInteractor {
	return &FilePermissionInteractorImpl{
		filePermissionRepo: repo,
		logger:             logger,
	}
}

func (f *FilePermissionInteractorImpl) GrantFilePermission(ctx context.Context, input *presenters.GrantFilePermissionInput) (*presenters.GrantFilePermissionOutput, error) {
	fileID, err := primitive.ObjectIDFromHex(input.FileID)
	if err != nil {
		return nil, err
	}
	userId, err := primitive.ObjectIDFromHex(input.UserID)
	if err != nil {
		return nil, err
	}
	filePermission := models.FilePermission{
		FileID:           fileID,
		UserID:           userId,
		PermissionType:   input.PermissionType,
		AccessSecureFile: input.AccessSecureFile,
	}
	res, err := f.filePermissionRepo.CreateFilePermission(ctx, &filePermission)
	if err != nil {
		return nil, err
	}
	f.logger.Info(ctx, map[string]interface{}{
		"file_id": fileID.Hex(),
		"user_id": userId.Hex(),
	}, "File granted success")
	return &presenters.GrantFilePermissionOutput{
		FileID:           res.FileID.Hex(),
		UserID:           res.UserID.Hex(),
		PermissionType:   input.PermissionType,
		AccessSecureFile: input.AccessSecureFile,
	}, nil
}
