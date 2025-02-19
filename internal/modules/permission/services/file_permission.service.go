package services

import (
	"storage-app/internal/modules/permission/repositories"
)

type PermissionService struct {
	Repo *repositories.PermissionRepository
}

func NewPermissionService() *PermissionService {
	return &PermissionService{
		Repo: repositories.NewPermissionRepository(),
	}
}

func (ps *PermissionService) GetAllPermissions() ([]map[string]interface{}, error) {
	return ps.Repo.FetchPermissions()
}

func (ps *PermissionService) GrantPermission(fileID, userID, permissionType string, accessSecure bool) error {
	// Business logic, validation, etc.
	if fileID == "" || userID == "" {
		return &InvalidInputError{"FileID and UserID are required"}
	}
	return ps.Repo.InsertPermission(fileID, userID, permissionType, accessSecure)
}

// Custom error handling
type InvalidInputError struct {
	Message string
}

func (e *InvalidInputError) Error() string {
	return e.Message
}
