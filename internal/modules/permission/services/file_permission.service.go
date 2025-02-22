package services

import (
	"strconv"

	"github.com/baothaihcmut/Storage-app/internal/modules/permission/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PermissionService struct {
	Repo *repositories.PermissionRepository
}

func NewPermissionService(db *mongo.Database) *PermissionService {
	return &PermissionService{
		Repo: repositories.NewPermissionRepository(db),
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

	fileObjectID, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
		return &InvalidInputError{"Invalid FileID format"}
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return &InvalidInputError{"Invalid UserID format"}
	}

	permTypeInt, err := strconv.Atoi(permissionType)
	if err != nil {
		return &InvalidInputError{"Invalid permissionType format"}
	}

	return ps.Repo.CreatePermission(fileObjectID, userObjectID, permTypeInt, accessSecure)
}

// Custom error handling
type InvalidInputError struct {
	Message string
}

func (e *InvalidInputError) Error() string {
	return e.Message
}
