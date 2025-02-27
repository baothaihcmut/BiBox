package controllers

import (
	"net/http"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_comment/interactors"

	"github.com/gin-gonic/gin"
)

type CommentController struct {
	Interactor *interactors.CommentInteractor
}

func NewCommentController(interactor *interactors.CommentInteractor) *CommentController {
	return &CommentController{
		Interactor: interactor,
	}
}

func (cc *CommentController) GetComments(c *gin.Context) {
	comments, err := cc.Interactor.GetAllComments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}
	c.JSON(http.StatusOK, comments)
}

func (cc *CommentController) AddComment(c *gin.Context) {
	var request struct {
		FileID  string `json:"file_id"`
		UserID  string `json:"user_id"`
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := cc.Interactor.AddComment(request.FileID, request.UserID, request.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment added"})
}
