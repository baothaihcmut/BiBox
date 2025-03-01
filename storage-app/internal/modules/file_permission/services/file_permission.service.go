package services

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/models"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreatePermssionArgs struct {
	FileID      primitive.ObjectID
	UserID      primitive.ObjectID
	Permssion   enums.FilePermissionType
	AcessSecure bool
	CanShare    bool
}

type PermissionService interface {
	CheckPermission(ctx context.Context, fileID, userId primitive.ObjectID, permssion enums.FilePermissionType) (bool, error)
	CreatePermssion(ctx context.Context, args CreatePermssionArgs) error
}

type PermissionServiceImpl struct {
	repo repositories.FilePermissionRepository
}

func NewPermissionService(repo repositories.FilePermissionRepository) PermissionService {
	return &PermissionServiceImpl{
		repo: repo,
	}
}

func (ps *PermissionServiceImpl) CreatePermssion(ctx context.Context, args CreatePermssionArgs) error {
	filePermssion := models.NewFilePermission(args.FileID, args.UserID, args.Permssion, args.AcessSecure, args.CanShare)
	return ps.repo.CreateFilePermission(ctx, filePermssion)
}

func (ps *PermissionServiceImpl) CheckPermission(ctx context.Context, fileID, userId primitive.ObjectID, permssion enums.FilePermissionType) (bool, error) {
	filePermssion, err := ps.repo.GetFilePermission(ctx, fileID, userId, repositories.FilterPermssionType{
		Option: repositories.PermssionGreaterThan,
		Value:  []enums.FilePermissionType{permssion},
	})
	if err != nil {
		return false, err
	}
	return filePermssion != nil, nil
}
