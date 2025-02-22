package interactors

import (
	"github.com/baothaihcmut/Storage-app/internal/modules/permission/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PermissionInteractor struct {
	Repo *repositories.PermissionRepository
}

func NewPermissionInteractor(db *mongo.Database) *PermissionInteractor {
	return &PermissionInteractor{
		Repo: repositories.NewPermissionRepository(db),
	}
}

func (pi *PermissionInteractor) GetAllPermissions() ([]map[string]interface{}, error) {
	return pi.Repo.FetchPermissions()
}

func (pi *PermissionInteractor) GrantPermission(fileID primitive.ObjectID, userID primitive.ObjectID, permissionType int, accessSecure bool) error {
	return pi.Repo.CreatePermission(fileID, userID, permissionType, accessSecure)
}
