package interactors

import (
	"storage-app/internal/modules/permission/repositories"
)

type PermissionInteractor struct {
	Repo *repositories.PermissionRepository
}

func NewPermissionInteractor() *PermissionInteractor {
	return &PermissionInteractor{
		Repo: repositories.NewPermissionRepository(),
	}
}

func (pi *PermissionInteractor) GetAllPermissions() ([]map[string]interface{}, error) {
	return pi.Repo.FetchPermissions()
}

func (pi *PermissionInteractor) GrantPermission(fileID, userID, permissionType string, accessSecure bool) error {
	return pi.Repo.InsertPermission(fileID, userID, permissionType, accessSecure)
}
