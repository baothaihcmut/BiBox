package services

import (
	"context"
	"time"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PermissionService struct {
	Repo *repositories.PermissionRepository
}

func NewPermissionService(repo *repositories.PermissionRepository) *PermissionService {
	return &PermissionService{
		Repo: repo,
	}
}

func (ps *PermissionService) UpdatePermission(fileID, userID primitive.ObjectID, permissionType int, accessSecure bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) //  Create context with timeout
	defer cancel()

	return ps.Repo.UpdatePermission(ctx, fileID, userID, permissionType, accessSecure)
}
