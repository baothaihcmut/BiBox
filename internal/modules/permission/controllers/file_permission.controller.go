package controllers

import (
	"net/http"
	"storage-app/internal/modules/permission/interactors"

	"github.com/gin-gonic/gin"
)

type PermissionController struct {
	Interactor *interactors.PermissionInteractor
}

func NewPermissionController() *PermissionController {
	return &PermissionController{
		Interactor: interactors.NewPermissionInteractor(),
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
		FileID         string `json:"file_id"`
		UserID         string `json:"user_id"`
		PermissionType string `json:"permission_type"`
		AccessSecure   bool   `json:"access_secure_file"`
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
