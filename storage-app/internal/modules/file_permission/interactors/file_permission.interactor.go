package interactors

import (
	"context"
	"time"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/constant"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/exception"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/models"
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

// Function to create file permission
func (pi *PermissionInteractor) CreateFilePermission(ctx context.Context, fileID string, canShare bool) error {
	// Get user context from token
	userContext, ok := ctx.Value(constant.UserContext).(*models.UserContext)
	if !ok {
		return exception.ErrUnauthorized
	}

	// Convert fileID to ObjectID
	fileObjectID, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
		return exception.ErrInvalidObjectId
	}

	// Fetch file details
	file, err := pi.Repo.GetFileByID(ctx, fileObjectID)
	if err != nil {
		return err
	}
	if file == nil {
		return exception.ErrFileNotFound
	}

	// Check if the user is the owner
	ownerUserID := file.OwnerUserID.Hex()
	if ownerUserID != userContext.Id {
		return exception.ErrPermissionDenied // User is not the owner
	}

	//Create file permission in DB
	err = pi.Repo.CreateFilePermission(ctx, fileObjectID, ownerUserID, canShare)
	if err != nil {
		return err
	}

	return nil
}
