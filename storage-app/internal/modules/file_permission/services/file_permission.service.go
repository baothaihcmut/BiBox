package services

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PermissionService interface {
	CheckPermission(ctx context.Context, fileID, userId primitive.ObjectID, permssion enums.FilePermissionType) (bool, error)
}

type PermissionServiceImpl struct {
	repo repositories.FilePermissionRepository
}

func NewPermissionService(repo repositories.FilePermissionRepository) PermissionService {
	return &PermissionServiceImpl{
		repo: repo,
	}
}

func (ps *PermissionServiceImpl) CheckPermission(ctx context.Context, fileID, userId primitive.ObjectID, permssion enums.FilePermissionType) (bool, error) {
	filePermssion, err := ps.repo.FindFilePermissionById(ctx, repositories.FilePermissionId{
		FileId: fileID,
		UserId: userId,
	})
	if err != nil {
		return false, err
	}
	if filePermssion == nil {
		return false, nil
	}
	return filePermssion.FilePermissionType >= permssion, nil

}
