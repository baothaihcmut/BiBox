package controllers

import (
	"net/http"

	"github.com/baothaihcmut/Storage-app/internal/modules/permission/interactors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PermissionController struct {
	Interactor *interactors.PermissionInteractor
}

func NewPermissionController(db *mongo.Database) *PermissionController {
	return &PermissionController{
		Interactor: interactors.NewPermissionInteractor(db),
	}
}
func (pc *PermissionController) GetPermissions(c *gin.Context) {
	permissions, err := pc.Interactor.GetAllPermissions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch permissions"})
		return
	}
	c.JSON(http.StatusOK, permissions)
}

func (pc *PermissionController) GrantPermission(c *gin.Context) {
	var request struct {
		FileID         primitive.ObjectID `json:"file_id"`
		UserID         primitive.ObjectID `json:"user_id"`
		PermissionType int                `json:"permission_type"`
		AccessSecure   bool               `json:"access_secure_file"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := pc.Interactor.GrantPermission(request.FileID, request.UserID, request.PermissionType, request.AccessSecure)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to grant permission"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permission granted"})
}
