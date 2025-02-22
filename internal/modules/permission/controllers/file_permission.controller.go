package controllers

import (
	"net/http"

	"github.com/baothaihcmut/Storage-app/internal/modules/permission/interactors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PermissionController struct {
	Interactor *interactors.PermissionInteractor
}

func NewPermissionController(interactor *interactors.PermissionInteractor) *PermissionController {
	return &PermissionController{
		Interactor: interactor,
	}
}

// update permission
func (pc *PermissionController) UpdatePermission(c *gin.Context) {
	var request struct {
		FileID         string `json:"file_id"`
		UserID         string `json:"user_id"`
		PermissionType int    `json:"permission_type"`
		AccessSecure   bool   `json:"access_secure_file"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// FileID and UserID from string to ObjectID
	fileObjectID, err := primitive.ObjectIDFromHex(request.FileID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file_id format"})
		return
	}

	userObjectID, err := primitive.ObjectIDFromHex(request.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id format"})
		return
	}

	// call the interactor
	err = pc.Interactor.UpdatePermission(fileObjectID, userObjectID, request.PermissionType, request.AccessSecure)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update permission"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permission updated successfully"})
}
