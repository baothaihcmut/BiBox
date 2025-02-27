package interactors

import (
	"context"
	"time"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PermissionInteractor struct {
	Repo *repositories.PermissionRepository
}

func NewPermissionInteractor(repo *repositories.PermissionRepository) *PermissionInteractor {
	return &PermissionInteractor{
		Repo: repo,
	}
}

// Update Permission
func (pi *PermissionInteractor) UpdatePermission(fileID primitive.ObjectID, userID primitive.ObjectID, permissionType int, accessSecure bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Call to repo to update Permission
	return pi.Repo.UpdatePermission(ctx, fileID, userID, permissionType, accessSecure)
}
